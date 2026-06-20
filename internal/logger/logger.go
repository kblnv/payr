package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
	FatalLevel Level = "FATAL"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format(time.RFC3339)
	log.Printf("[%s] %s %s", level, timestamp, msg)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
	os.Exit(1)
}
