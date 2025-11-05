package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meysam81/parse-dmarc/internal/api"
	"github.com/meysam81/parse-dmarc/internal/config"
	"github.com/meysam81/parse-dmarc/internal/imap"
	"github.com/meysam81/parse-dmarc/internal/parser"
	"github.com/meysam81/parse-dmarc/internal/storage"
)

func main() {
	var (
		configPath  = flag.String("config", "config.json", "Path to configuration file")
		genConfig   = flag.Bool("gen-config", false, "Generate sample configuration file")
		fetchOnce   = flag.Bool("fetch-once", false, "Fetch reports once and exit")
		serveOnly   = flag.Bool("serve-only", false, "Only serve the dashboard without fetching")
		fetchInterval = flag.Int("fetch-interval", 300, "Interval in seconds between fetch operations")
	)
	flag.Parse()

	// Generate sample config if requested
	if *genConfig {
		if err := config.GenerateSample(*configPath); err != nil {
			log.Fatalf("Failed to generate config: %v", err)
		}
		log.Printf("Sample configuration written to %s", *configPath)
		return
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize storage
	store, err := storage.NewStorage(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	// Start API server in background
	server := api.NewServer(store, cfg.Server.Host, cfg.Server.Port)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// If serve-only mode, just wait for signal
	if *serveOnly {
		log.Println("Running in serve-only mode")
		waitForSignal()
		return
	}

	// Fetch reports
	if *fetchOnce {
		if err := fetchReports(cfg, store); err != nil {
			log.Fatalf("Failed to fetch reports: %v", err)
		}
		log.Println("Fetch complete")
		return
	}

	// Continuous fetching
	log.Printf("Starting continuous fetch mode (interval: %d seconds)", *fetchInterval)

	// Initial fetch
	if err := fetchReports(cfg, store); err != nil {
		log.Printf("Initial fetch failed: %v", err)
	}

	// Set up ticker for periodic fetching
	ticker := time.NewTicker(time.Duration(*fetchInterval) * time.Second)
	defer ticker.Stop()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			if err := fetchReports(cfg, store); err != nil {
				log.Printf("Fetch failed: %v", err)
			}
		case <-sigChan:
			log.Println("Shutting down...")
			return
		}
	}
}

func fetchReports(cfg *config.Config, store *storage.Storage) error {
	log.Println("Fetching DMARC reports...")

	// Create IMAP client
	client := imap.NewClient(&cfg.IMAP)
	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Disconnect()

	// Fetch reports
	reports, err := client.FetchDMARCReports()
	if err != nil {
		return err
	}

	if len(reports) == 0 {
		log.Println("No new reports found")
		return nil
	}

	log.Printf("Processing %d reports...", len(reports))

	// Process each report
	processed := 0
	for _, report := range reports {
		for _, attachment := range report.Attachments {
			feedback, err := parser.ParseReport(attachment.Data)
			if err != nil {
				log.Printf("Failed to parse %s: %v", attachment.Filename, err)
				continue
			}

			if err := store.SaveReport(feedback); err != nil {
				log.Printf("Failed to save report %s: %v", feedback.ReportMetadata.ReportID, err)
				continue
			}

			log.Printf("Saved report: %s from %s (domain: %s, messages: %d)",
				feedback.ReportMetadata.ReportID,
				feedback.ReportMetadata.OrgName,
				feedback.PolicyPublished.Domain,
				feedback.GetTotalMessages())
			processed++
		}
	}

	log.Printf("Successfully processed %d reports", processed)
	return nil
}

func waitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
