package migration

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// AuditResult holds the pre-migration audit for a single database.
type AuditResult struct {
	Database string            `json:"database"`
	Tables   []string          `json:"tables,omitempty"`
	Counts   map[string]int64  `json:"counts,omitempty"`
	Issues   map[string]string `json:"issues,omitempty"`
}

// DedupResult holds the outcome of a deduplication operation.
type DedupResult struct {
	Table          string `json:"table"`
	DuplicatesSeen int64  `json:"duplicates_seen"`
	Removed        int64  `json:"removed"`
	Kept           int64  `json:"kept"`
}

// ValidationCheck holds a single validation result.
type ValidationCheck struct {
	Name     string      `json:"name"`
	Status   string      `json:"status"` // "PASS" or "FAIL"
	Expected interface{} `json:"expected,omitempty"`
	Actual   interface{} `json:"actual,omitempty"`
	Details  string      `json:"details,omitempty"`
}

// Report holds the complete migration report.
type Report struct {
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt time.Time         `json:"completed_at"`
	Duration    string            `json:"duration"`
	DryRun      bool              `json:"dry_run"`
	Audit       []AuditResult     `json:"audit,omitempty"`
	DataCopied  map[string]int64  `json:"data_copied,omitempty"`
	DataSkipped map[string]int64  `json:"data_skipped,omitempty"`
	Dedup       []DedupResult     `json:"dedup,omitempty"`
	Validation  []ValidationCheck `json:"validation,omitempty"`
	Overall     string            `json:"overall"` // "PASS" or "FAIL"
}

// NewReport creates a new report with the start time set.
func NewReport(dryRun bool) *Report {
	return &Report{
		StartedAt:   time.Now(),
		DryRun:      dryRun,
		DataCopied:  make(map[string]int64),
		DataSkipped: make(map[string]int64),
	}
}

// Finish marks the report as complete and calculates duration.
func (r *Report) Finish() {
	r.CompletedAt = time.Now()
	r.Duration = r.CompletedAt.Sub(r.StartedAt).Round(time.Millisecond).String()

	r.Overall = "PASS"
	for _, v := range r.Validation {
		if v.Status == "FAIL" {
			r.Overall = "FAIL"
			break
		}
	}
}

// PrintJSON writes the report as formatted JSON to w.
func (r *Report) PrintJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// PrintSummary writes a human-readable summary to w.
func (r *Report) PrintSummary(w io.Writer) {
	fmt.Fprintf(w, "\n=== Migration Report ===\n")
	fmt.Fprintf(w, "Duration: %s\n", r.Duration)
	fmt.Fprintf(w, "Dry Run:  %v\n", r.DryRun)

	if len(r.DataCopied) > 0 {
		fmt.Fprintf(w, "\nData Copied:\n")
		for table, count := range r.DataCopied {
			fmt.Fprintf(w, "  %-40s %d rows\n", table, count)
		}
	}

	if len(r.DataSkipped) > 0 {
		fmt.Fprintf(w, "\nData Skipped (conflicts):\n")
		for table, count := range r.DataSkipped {
			fmt.Fprintf(w, "  %-40s %d rows\n", table, count)
		}
	}

	if len(r.Dedup) > 0 {
		fmt.Fprintf(w, "\nDeduplication:\n")
		for _, d := range r.Dedup {
			fmt.Fprintf(w, "  %-40s %d duplicates found, %d removed\n", d.Table, d.DuplicatesSeen, d.Removed)
		}
	}

	if len(r.Validation) > 0 {
		fmt.Fprintf(w, "\nValidation:\n")
		for _, v := range r.Validation {
			fmt.Fprintf(w, "  [%s] %s", v.Status, v.Name)
			if v.Details != "" {
				fmt.Fprintf(w, " — %s", v.Details)
			}
			fmt.Fprintln(w)
		}
	}

	fmt.Fprintf(w, "\nOverall: %s\n", r.Overall)
}
