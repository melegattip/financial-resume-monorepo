package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Setup initializes the global zerolog logger with the given log level.
// In development mode, it uses a human-readable console writer.
// In production, it outputs structured JSON.
func Setup(level string, env string) zerolog.Logger {
	var output io.Writer

	if env == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		output = os.Stdout
	}

	lvl := parseLevel(level)

	logger := zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
