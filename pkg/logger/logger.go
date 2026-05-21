package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger is a simple logger wrapper
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger with the given level
func New(levelStr string) *Logger {
	level := INFO
	switch strings.ToLower(levelStr) {
	case "debug":
		level = DEBUG
	case "info":
		level = INFO
	case "warn":
		level = WARN
	case "error":
		level = ERROR
	}

	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Printf logs at info level
func (l *Logger) Printf(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[%s] INFO: %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, v...))
	}
}

// Debugf logs at debug level
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[%s] DEBUG: %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, v...))
	}
}

// Warnf logs at warn level
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[%s] WARN: %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, v...))
	}
}

// Errorf logs at error level
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[%s] ERROR: %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, v...))
	}
}
