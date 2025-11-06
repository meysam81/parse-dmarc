package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/meysam81/parse-dmarc/internal/api"
	"github.com/meysam81/parse-dmarc/internal/config"
	"github.com/meysam81/parse-dmarc/internal/imap"
	"github.com/meysam81/parse-dmarc/internal/parser"
	"github.com/meysam81/parse-dmarc/internal/storage"
	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	cmd := &cli.Command{
		Name:                  "parse-dmarc",
		Usage:                 "Parse and analyze DMARC reports from IMAP mailbox",
		Version:               version,
		EnableShellCompletion: true,
		Suggest:               true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to configuration file",
				Value:   "config.json",
				Sources: cli.EnvVars("PARSE_DMARC_CONFIG"),
			},
			&cli.BoolFlag{
				Name:    "gen-config",
				Usage:   "Generate sample configuration file",
				Sources: cli.EnvVars("PARSE_DMARC_GEN_CONFIG"),
			},
			&cli.BoolFlag{
				Name:    "fetch-once",
				Usage:   "Fetch reports once and exit",
				Sources: cli.EnvVars("PARSE_DMARC_FETCH_ONCE"),
			},
			&cli.BoolFlag{
				Name:    "serve-only",
				Usage:   "Only serve the dashboard without fetching",
				Sources: cli.EnvVars("PARSE_DMARC_SERVE_ONLY"),
			},
			&cli.IntFlag{
				Name:    "fetch-interval",
				Usage:   "Interval in seconds between fetch operations",
				Value:   300,
				Sources: cli.EnvVars("PARSE_DMARC_FETCH_INTERVAL"),
			},
		},
		Action: run,
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version information",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Version:    %s\n", version)
					fmt.Printf("Commit:     %s\n", commit)
					fmt.Printf("Build Date: %s\n", date)
					fmt.Printf("Built By:   %s\n", builtBy)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	configPath := cmd.String("config")
	genConfig := cmd.Bool("gen-config")
	fetchOnce := cmd.Bool("fetch-once")
	serveOnly := cmd.Bool("serve-only")
	fetchInterval := cmd.Int("fetch-interval")

	if genConfig {
		if err := config.GenerateSample(configPath); err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
		log.Printf("Sample configuration written to %s", configPath)
		return nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer func() { _ = store.Close() }()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := api.NewServer(store, cfg.Server.Host, cfg.Server.Port)
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- server.Start(ctx)
	}()

	if serveOnly {
		log.Println("Running in serve-only mode")
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
		return nil
	}

	if fetchOnce {
		if err := fetchReports(cfg, store); err != nil {
			return fmt.Errorf("failed to fetch reports: %w", err)
		}
		log.Println("Fetch complete")
		return nil
	}

	log.Printf("Starting continuous fetch mode (interval: %d seconds)", fetchInterval)

	if err := fetchReports(cfg, store); err != nil {
		log.Printf("Initial fetch failed: %v", err)
	}

	ticker := time.NewTicker(time.Duration(fetchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fetchReports(cfg, store); err != nil {
				log.Printf("Fetch failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("Shutting down...")
			return nil
		case err := <-serverErrChan:
			if err != nil {
				return fmt.Errorf("server error: %w", err)
			}
		}
	}
}

func fetchReports(cfg *config.Config, store *storage.Storage) error {
	log.Println("Fetching DMARC reports...")

	// Get all IMAP configurations (supports both single and multiple inboxes)
	imapConfigs := cfg.GetIMAPConfigs()
	if len(imapConfigs) == 0 {
		return fmt.Errorf("no IMAP configuration found")
	}

	log.Printf("Fetching from %d inbox(es)...", len(imapConfigs))

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error
	totalProcessed := 0

	// Fetch from each IMAP inbox concurrently
	for i := range imapConfigs {
		wg.Add(1)
		go func(imapCfg *config.IMAPConfig, index int) {
			defer wg.Done()

			mailboxName := fmt.Sprintf("%s@%s:%s", imapCfg.Username, imapCfg.Host, imapCfg.Mailbox)
			log.Printf("[Inbox %d/%d] Connecting to %s", index+1, len(imapConfigs), mailboxName)

			// Create IMAP client for this inbox
			client := imap.NewClient(imapCfg)
			if err := client.Connect(); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to connect to %s: %w", mailboxName, err))
				mu.Unlock()
				return
			}
			defer func() { _ = client.Disconnect() }()

			// Fetch reports from this inbox
			reports, err := client.FetchDMARCReports()
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to fetch from %s: %w", mailboxName, err))
				mu.Unlock()
				return
			}

			if len(reports) == 0 {
				log.Printf("[Inbox %d/%d] No new reports found in %s", index+1, len(imapConfigs), mailboxName)
				return
			}

			log.Printf("[Inbox %d/%d] Processing %d report(s) from %s", index+1, len(imapConfigs), len(reports), mailboxName)

			// Process each report
			processed := 0
			for _, report := range reports {
				for _, attachment := range report.Attachments {
					feedback, err := parser.ParseReport(attachment.Data)
					if err != nil {
						log.Printf("[Inbox %d/%d] Failed to parse %s: %v", index+1, len(imapConfigs), attachment.Filename, err)
						continue
					}

					if err := store.SaveReport(feedback); err != nil {
						log.Printf("[Inbox %d/%d] Failed to save report %s: %v", index+1, len(imapConfigs), feedback.ReportMetadata.ReportID, err)
						continue
					}

					log.Printf("[Inbox %d/%d] Saved report: %s from %s (domain: %s, messages: %d)",
						index+1, len(imapConfigs),
						feedback.ReportMetadata.ReportID,
						feedback.ReportMetadata.OrgName,
						feedback.PolicyPublished.Domain,
						feedback.GetTotalMessages())
					processed++
				}
			}

			mu.Lock()
			totalProcessed += processed
			mu.Unlock()

			log.Printf("[Inbox %d/%d] Successfully processed %d report(s) from %s", index+1, len(imapConfigs), processed, mailboxName)
		}(&imapConfigs[i], i)
	}

	// Wait for all fetches to complete
	wg.Wait()

	// Report results
	if len(errors) > 0 {
		log.Printf("Completed with %d error(s):", len(errors))
		for _, err := range errors {
			log.Printf("  - %v", err)
		}
	}

	log.Printf("Successfully processed %d total report(s) across all inboxes", totalProcessed)

	// Return error only if all fetches failed
	if len(errors) > 0 && len(errors) == len(imapConfigs) {
		return fmt.Errorf("all inbox fetches failed")
	}

	return nil
}
