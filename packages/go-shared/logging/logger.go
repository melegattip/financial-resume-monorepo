package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the shared logger instance
var Logger *logrus.Logger

func init() {
	Logger = NewLogger()
}

// NewLogger creates a new configured logger
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	// Set log level from environment
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}

// WithService adds service name to log context
func WithService(serviceName string) *logrus.Entry {
	return Logger.WithField("service", serviceName)
}

// WithRequestID adds request ID to log context
func WithRequestID(requestID string) *logrus.Entry {
	return Logger.WithField("request_id", requestID)
}
