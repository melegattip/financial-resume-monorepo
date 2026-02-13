package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetup_Development(t *testing.T) {
	logger := Setup("debug", "development")
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
}

func TestSetup_Production(t *testing.T) {
	logger := Setup("info", "production")
	assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"unknown", zerolog.InfoLevel},
		{"", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, parseLevel(tt.input))
		})
	}
}
