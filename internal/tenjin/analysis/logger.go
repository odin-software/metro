package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// MaxLogFileSize is 5MB in bytes
	MaxLogFileSize = 5 * 1024 * 1024
)

// MetricsLogger handles writing metrics to rotating log files
type MetricsLogger struct {
	directory   string
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
}

// NewMetricsLogger creates a new metrics logger
func NewMetricsLogger(directory string) (*MetricsLogger, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metrics directory: %w", err)
	}

	logger := &MetricsLogger{
		directory: directory,
	}

	// Open initial log file
	if err := logger.rotate(); err != nil {
		return nil, err
	}

	return logger, nil
}

// Log writes metrics output to the current log file
func (l *MetricsLogger) Log(output string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if we need to rotate
	if l.currentSize >= MaxLogFileSize {
		if err := l.rotate(); err != nil {
			return err
		}
	}

	// Write to file
	n, err := l.currentFile.WriteString(output)
	if err != nil {
		return fmt.Errorf("failed to write metrics: %w", err)
	}

	l.currentSize += int64(n)
	return nil
}

// rotate closes the current file and opens a new one
func (l *MetricsLogger) rotate() error {
	// Close current file if it exists
	if l.currentFile != nil {
		if err := l.currentFile.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	// Generate new filename with timestamp
	filename := fmt.Sprintf("tenjin-metrics-%s.log",
		time.Now().Format("2006-01-02-150405"))
	filepath := filepath.Join(l.directory, filename)

	// Open new file
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.currentFile = file
	l.currentSize = 0

	// Write header
	header := fmt.Sprintf("=== Tenjin Metrics Log Started: %s ===\n\n",
		time.Now().Format("2006-01-02 15:04:05"))
	n, _ := l.currentFile.WriteString(header)
	l.currentSize += int64(n)

	return nil
}

// Close closes the current log file
func (l *MetricsLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentFile != nil {
		return l.currentFile.Close()
	}
	return nil
}
