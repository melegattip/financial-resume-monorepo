package migration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplitName tests the name splitting helper used during user migration.
func TestSplitName(t *testing.T) {
	tests := []struct {
		input     string
		firstName string
		lastName  string
	}{
		{"John Doe", "John", "Doe"},
		{"Jane", "Jane", ""},
		{"John Michael Doe", "John", "Michael Doe"},
		{"", "", ""},
		{"  Alice  ", "Alice", ""},
		{"  Bob Smith  ", "Bob", "Smith"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			first, last := splitName(tt.input)
			assert.Equal(t, tt.firstName, first)
			assert.Equal(t, tt.lastName, last)
		})
	}
}

// TestAuditResult verifies AuditResult struct initialization.
func TestAuditResult(t *testing.T) {
	result := AuditResult{
		Database: "test_db",
		Counts:   map[string]int64{"users": 10},
		Issues:   map[string]string{},
	}

	assert.Equal(t, "test_db", result.Database)
	assert.Equal(t, int64(10), result.Counts["users"])
	assert.Empty(t, result.Issues)
}

// TestDedupResult verifies DedupResult struct.
func TestDedupResult(t *testing.T) {
	result := DedupResult{
		Table:          "user_gamification",
		DuplicatesSeen: 5,
		Removed:        4,
		Kept:           1,
	}

	assert.Equal(t, "user_gamification", result.Table)
	assert.Equal(t, int64(5), result.DuplicatesSeen)
	assert.Equal(t, int64(4), result.Removed)
}

// TestValidationCheck verifies ValidationCheck struct.
func TestValidationCheck(t *testing.T) {
	check := ValidationCheck{
		Name:     "user_count",
		Status:   "PASS",
		Expected: int64(10),
		Actual:   int64(10),
	}

	assert.Equal(t, "PASS", check.Status)
	assert.Equal(t, int64(10), check.Expected)
}

// TestReportFinish verifies report completion logic.
func TestReportFinish_AllPass(t *testing.T) {
	report := NewReport(false)
	report.Validation = []ValidationCheck{
		{Name: "check1", Status: "PASS"},
		{Name: "check2", Status: "PASS"},
	}
	report.Finish()

	assert.Equal(t, "PASS", report.Overall)
	assert.NotEmpty(t, report.Duration)
	assert.False(t, report.CompletedAt.IsZero())
}

func TestReportFinish_AnyFail(t *testing.T) {
	report := NewReport(false)
	report.Validation = []ValidationCheck{
		{Name: "check1", Status: "PASS"},
		{Name: "check2", Status: "FAIL"},
		{Name: "check3", Status: "PASS"},
	}
	report.Finish()

	assert.Equal(t, "FAIL", report.Overall)
}

func TestReportFinish_DryRun(t *testing.T) {
	report := NewReport(true)
	report.Finish()

	assert.True(t, report.DryRun)
	assert.Equal(t, "PASS", report.Overall) // No validation checks = PASS
}

// TestNewReport verifies initial report state.
func TestNewReport(t *testing.T) {
	report := NewReport(true)

	assert.True(t, report.DryRun)
	assert.NotNil(t, report.DataCopied)
	assert.NotNil(t, report.DataSkipped)
	assert.False(t, report.StartedAt.IsZero())
}
