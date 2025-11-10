package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger provides structured JSON logging
type Logger struct {
	output *os.File
}

// New creates a new logger instance that writes to stdout
func New() *Logger {
	return &Logger{
		output: os.Stdout,
	}
}

// log writes a log entry in JSON format
func (l *Logger) log(level LogLevel, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple format if JSON marshaling fails
		fmt.Fprintf(l.output, "{\"error\":\"failed to marshal log entry\",\"message\":\"%s\"}\n", message)
		return
	}

	fmt.Fprintln(l.output, string(jsonData))
}

// Info logs an info level message
func (l *Logger) Info(message string) {
	l.log(LevelInfo, message, nil)
}

// Infof logs an info level message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// InfoWithFields logs an info level message with additional fields
func (l *Logger) InfoWithFields(message string, fields map[string]interface{}) {
	l.log(LevelInfo, message, fields)
}

// Error logs an error level message
func (l *Logger) Error(message string) {
	l.log(LevelError, message, nil)
}

// Errorf logs an error level message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// ErrorWithFields logs an error level message with additional fields
func (l *Logger) ErrorWithFields(message string, fields map[string]interface{}) {
	l.log(LevelError, message, fields)
}

// Warn logs a warning level message
func (l *Logger) Warn(message string) {
	l.log(LevelWarn, message, nil)
}

// Warnf logs a warning level message with formatting
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// WarnWithFields logs a warning level message with additional fields
func (l *Logger) WarnWithFields(message string, fields map[string]interface{}) {
	l.log(LevelWarn, message, fields)
}

// Debug logs a debug level message
func (l *Logger) Debug(message string) {
	l.log(LevelDebug, message, nil)
}

// Debugf logs a debug level message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// DebugWithFields logs a debug level message with additional fields
func (l *Logger) DebugWithFields(message string, fields map[string]interface{}) {
	l.log(LevelDebug, message, fields)
}

// Default logger instance for convenience
var defaultLogger = New()

// Info logs an info level message using the default logger
func Info(message string) {
	defaultLogger.Info(message)
}

// Infof logs an info level message with formatting using the default logger
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Error logs an error level message using the default logger
func Error(message string) {
	defaultLogger.Error(message)
}

// Errorf logs an error level message with formatting using the default logger
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Warn logs a warning level message using the default logger
func Warn(message string) {
	defaultLogger.Warn(message)
}

// Warnf logs a warning level message with formatting using the default logger
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Debug logs a debug level message using the default logger
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Debugf logs a debug level message with formatting using the default logger
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// ErrorWithFields logs an error level message with additional fields using the default logger
func ErrorWithFields(message string, fields map[string]interface{}) {
	defaultLogger.ErrorWithFields(message, fields)
}

// InfoWithFields logs an info level message with additional fields using the default logger
func InfoWithFields(message string, fields map[string]interface{}) {
	defaultLogger.InfoWithFields(message, fields)
}

// WarnWithFields logs a warning level message with additional fields using the default logger
func WarnWithFields(message string, fields map[string]interface{}) {
	defaultLogger.WarnWithFields(message, fields)
}

// DebugWithFields logs a debug level message with additional fields using the default logger
func DebugWithFields(message string, fields map[string]interface{}) {
	defaultLogger.DebugWithFields(message, fields)
}
