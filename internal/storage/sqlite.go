package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/meysam81/parse-dmarc/internal/parser"
)

// Storage handles database operations
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new storage instance
func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &Storage{db: db}
	if err := storage.init(); err != nil {
		return nil, err
	}

	return storage, nil
}

// init initializes database schema
func (s *Storage) init() error {
	schema := `
	CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_id TEXT UNIQUE NOT NULL,
		org_name TEXT NOT NULL,
		email TEXT,
		domain TEXT NOT NULL,
		date_begin INTEGER NOT NULL,
		date_end INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		policy_p TEXT,
		policy_sp TEXT,
		policy_pct INTEGER,
		total_messages INTEGER,
		compliant_messages INTEGER,
		raw_report TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		report_id INTEGER NOT NULL,
		source_ip TEXT NOT NULL,
		count INTEGER NOT NULL,
		disposition TEXT,
		dkim_result TEXT,
		spf_result TEXT,
		header_from TEXT,
		envelope_from TEXT,
		dkim_domains TEXT,
		spf_domains TEXT,
		FOREIGN KEY (report_id) REFERENCES reports(id)
	);

	CREATE INDEX IF NOT EXISTS idx_reports_date_begin ON reports(date_begin);
	CREATE INDEX IF NOT EXISTS idx_reports_domain ON reports(domain);
	CREATE INDEX IF NOT EXISTS idx_records_report_id ON records(report_id);
	CREATE INDEX IF NOT EXISTS idx_records_source_ip ON records(source_ip);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveReport saves a DMARC report to the database
func (s *Storage) SaveReport(feedback *parser.Feedback) error {
	// Serialize the full report
	rawReport, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert report metadata
	result, err := tx.Exec(`
		INSERT OR IGNORE INTO reports (
			report_id, org_name, email, domain,
			date_begin, date_end, created_at,
			policy_p, policy_sp, policy_pct,
			total_messages, compliant_messages,
			raw_report
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		feedback.ReportMetadata.ReportID,
		feedback.ReportMetadata.OrgName,
		feedback.ReportMetadata.Email,
		feedback.PolicyPublished.Domain,
		feedback.ReportMetadata.DateRange.Begin,
		feedback.ReportMetadata.DateRange.End,
		time.Now().Unix(),
		feedback.PolicyPublished.P,
		feedback.PolicyPublished.SP,
		feedback.PolicyPublished.PCT,
		feedback.GetTotalMessages(),
		feedback.GetDMARCCompliantCount(),
		rawReport,
	)

	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	reportID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// If report already exists (rowsAffected == 0), skip records
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil // Report already exists
	}

	// Insert records
	for _, record := range feedback.Records {
		dkimDomains, _ := json.Marshal(record.AuthResults.DKIM)
		spfDomains, _ := json.Marshal(record.AuthResults.SPF)

		_, err := tx.Exec(`
			INSERT INTO records (
				report_id, source_ip, count,
				disposition, dkim_result, spf_result,
				header_from, envelope_from,
				dkim_domains, spf_domains
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			reportID,
			record.Row.SourceIP,
			record.Row.Count,
			record.Row.PolicyEvaluated.Disposition,
			record.Row.PolicyEvaluated.DKIM,
			record.Row.PolicyEvaluated.SPF,
			record.Identifiers.HeaderFrom,
			record.Identifiers.EnvelopeFrom,
			dkimDomains,
			spfDomains,
		)

		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	return tx.Commit()
}

// ReportSummary represents a summary of a report
type ReportSummary struct {
	ID                 int64  `json:"id"`
	ReportID           string `json:"report_id"`
	OrgName            string `json:"org_name"`
	Domain             string `json:"domain"`
	DateBegin          int64  `json:"date_begin"`
	DateEnd            int64  `json:"date_end"`
	TotalMessages      int    `json:"total_messages"`
	CompliantMessages  int    `json:"compliant_messages"`
	ComplianceRate     float64 `json:"compliance_rate"`
	PolicyP            string `json:"policy_p"`
}

// GetReports retrieves all reports
func (s *Storage) GetReports(limit, offset int) ([]ReportSummary, error) {
	rows, err := s.db.Query(`
		SELECT id, report_id, org_name, domain,
		       date_begin, date_end,
		       total_messages, compliant_messages,
		       policy_p
		FROM reports
		ORDER BY date_begin DESC
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []ReportSummary
	for rows.Next() {
		var r ReportSummary
		err := rows.Scan(
			&r.ID, &r.ReportID, &r.OrgName, &r.Domain,
			&r.DateBegin, &r.DateEnd,
			&r.TotalMessages, &r.CompliantMessages,
			&r.PolicyP,
		)
		if err != nil {
			return nil, err
		}

		if r.TotalMessages > 0 {
			r.ComplianceRate = float64(r.CompliantMessages) / float64(r.TotalMessages) * 100
		}

		reports = append(reports, r)
	}

	return reports, nil
}

// GetReportByID retrieves a full report by ID
func (s *Storage) GetReportByID(id int64) (*parser.Feedback, error) {
	var rawReport string
	err := s.db.QueryRow("SELECT raw_report FROM reports WHERE id = ?", id).Scan(&rawReport)
	if err != nil {
		return nil, err
	}

	var feedback parser.Feedback
	if err := json.Unmarshal([]byte(rawReport), &feedback); err != nil {
		return nil, err
	}

	return &feedback, nil
}

// Statistics holds dashboard statistics
type Statistics struct {
	TotalReports       int     `json:"total_reports"`
	TotalMessages      int     `json:"total_messages"`
	CompliantMessages  int     `json:"compliant_messages"`
	ComplianceRate     float64 `json:"compliance_rate"`
	UniqueSourceIPs    int     `json:"unique_source_ips"`
	UniqueDomains      int     `json:"unique_domains"`
}

// GetStatistics retrieves dashboard statistics
func (s *Storage) GetStatistics() (*Statistics, error) {
	var stats Statistics

	// Get total reports and messages
	err := s.db.QueryRow(`
		SELECT
			COUNT(*) as total_reports,
			COALESCE(SUM(total_messages), 0) as total_messages,
			COALESCE(SUM(compliant_messages), 0) as compliant_messages
		FROM reports
	`).Scan(&stats.TotalReports, &stats.TotalMessages, &stats.CompliantMessages)

	if err != nil {
		return nil, err
	}

	if stats.TotalMessages > 0 {
		stats.ComplianceRate = float64(stats.CompliantMessages) / float64(stats.TotalMessages) * 100
	}

	// Get unique source IPs
	err = s.db.QueryRow("SELECT COUNT(DISTINCT source_ip) FROM records").Scan(&stats.UniqueSourceIPs)
	if err != nil {
		return nil, err
	}

	// Get unique domains
	err = s.db.QueryRow("SELECT COUNT(DISTINCT domain) FROM reports").Scan(&stats.UniqueDomains)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// TopSourceIP represents source IP statistics
type TopSourceIP struct {
	SourceIP string `json:"source_ip"`
	Count    int    `json:"count"`
	Pass     int    `json:"pass"`
	Fail     int    `json:"fail"`
}

// GetTopSourceIPs retrieves top source IPs
func (s *Storage) GetTopSourceIPs(limit int) ([]TopSourceIP, error) {
	rows, err := s.db.Query(`
		SELECT
			source_ip,
			SUM(count) as total_count,
			SUM(CASE WHEN (dkim_result = 'pass' OR spf_result = 'pass') THEN count ELSE 0 END) as pass_count,
			SUM(CASE WHEN (dkim_result != 'pass' AND spf_result != 'pass') THEN count ELSE 0 END) as fail_count
		FROM records
		GROUP BY source_ip
		ORDER BY total_count DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TopSourceIP
	for rows.Next() {
		var r TopSourceIP
		if err := rows.Scan(&r.SourceIP, &r.Count, &r.Pass, &r.Fail); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}
