package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus for consistent logging across the application
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger(level string) *Logger {
	logger := logrus.New()
	
	// Set log format
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	
	// Set output to stdout
	logger.SetOutput(os.Stdout)
	
	// Parse and set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)
	
	return &Logger{Logger: logger}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.Logger.Info(message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.Logger.Error(message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.Logger.Debug(message)
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	l.Logger.Warn(message)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string) {
	l.Logger.Fatal(message)
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
