package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Init initializes the global logger with structured logging
func Init() {
	// Set log format
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set output
	logrus.SetOutput(os.Stdout)

	// Add default fields
	logrus.AddHook(&DefaultFieldsHook{})
}

// DefaultFieldsHook adds default fields to all log entries
type DefaultFieldsHook struct{}

func (hook *DefaultFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *DefaultFieldsHook) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = "block-flow"
	return nil
}
