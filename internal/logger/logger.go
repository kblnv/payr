package logger

import (
	"fmt"
	"log"
	"os"
)

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
	FatalLevel Level = "FATAL"
)

type Logger struct {
	pkg string
}

func New() *Logger {
	return &Logger{pkg: ""}
}

func (l *Logger) WithPackage(pkg string) *Logger {
	l.pkg = pkg
	return l
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.pkg == "" {
		log.Printf("level=%-5s msg=%s", level, msg)
	} else {
		log.Printf("level=%s package=%s msg=%s", level, l.pkg, msg)
	}
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
