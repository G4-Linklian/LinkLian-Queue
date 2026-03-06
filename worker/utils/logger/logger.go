package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	LevelLog     LogLevel = "LOG"
	LevelError   LogLevel = "ERROR"
	LevelWarn    LogLevel = "WARN"
	LevelDebug   LogLevel = "DEBUG"
	LevelVerbose LogLevel = "VERBOSE"
)

// AppLogger is a structured logger similar to NestJS Winston AppLogger
type AppLogger struct {
	logger *log.Logger
}

// New creates a new AppLogger instance
func New() *AppLogger {
	return &AppLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// Default is a package-level logger instance
var Default = New()

func (l *AppLogger) write(level LogLevel, message string, from string, detail any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	output := fmt.Sprintf("%s %s [%s] %s", timestamp, level, from, message)

	if detail != nil {
		detailBytes, err := json.MarshalIndent(detail, "", "  ")
		if err == nil {
			output += " " + string(detailBytes)
		}
	}

	l.logger.Println(output)
}

// Log logs a general message
func (l *AppLogger) Log(message string, from string, detail ...any) {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	l.write(LevelLog, message, from, d)
}

// Error logs an error message
func (l *AppLogger) Error(message string, from string, detail ...any) {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	l.write(LevelError, message, from, d)
}

// Warn logs a warning message
func (l *AppLogger) Warn(message string, from string, detail ...any) {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	l.write(LevelWarn, message, from, d)
}

// Debug logs a debug message
func (l *AppLogger) Debug(message string, from string, detail ...any) {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	l.write(LevelDebug, message, from, d)
}

// Verbose logs a verbose message
func (l *AppLogger) Verbose(message string, from string, detail ...any) {
	var d any
	if len(detail) > 0 {
		d = detail[0]
	}
	l.write(LevelVerbose, message, from, d)
}

// --- Package-level functions using Default logger ---

func Log(message string, from string, detail ...any) {
	Default.Log(message, from, detail...)
}

func Error(message string, from string, detail ...any) {
	Default.Error(message, from, detail...)
}

func Warn(message string, from string, detail ...any) {
	Default.Warn(message, from, detail...)
}

func Debug(message string, from string, detail ...any) {
	Default.Debug(message, from, detail...)
}

func Verbose(message string, from string, detail ...any) {
	Default.Verbose(message, from, detail...)
}
