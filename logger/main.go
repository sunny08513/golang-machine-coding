package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	SetOutput(writer io.Writer)
}

type LogLevel int

const (
	InfoLevel LogLevel = iota
	WarnLevel
	ErrorLevel
)

type SimpleLogger struct {
	mu     sync.Mutex
	writer io.Writer
}

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{writer: os.Stdout}
}

func (l *SimpleLogger) log(level LogLevel, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	var levelStr string
	switch level {
	case InfoLevel:
		levelStr = "INFO"
	case WarnLevel:
		levelStr = "WARN"
	case ErrorLevel:
		levelStr = "ERROR"
	}

	logMsg := fmt.Sprintf("%s [%s] %s\n", timestamp, levelStr, msg)
	l.writer.Write([]byte(logMsg))
}

func (l *SimpleLogger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *SimpleLogger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *SimpleLogger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *SimpleLogger) SetOutput(writer io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = writer
}

func main() {
	log := NewSimpleLogger()

	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")

	// Redirect log output to a file
	file, err := os.Create("log.txt")
	if err != nil {
		log.Error("Failed to create log file")
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.Info("This message goes to the log file")
}
