package log

import (
	"fmt"
	"os"
	"sync"
	"time"

)

type RotatingLogger struct {
    file    *os.File
    path    string
    maxSize int64
    mu      sync.Mutex 
}

// NewRotatingLogger creates thread-safe logger
func NewRotatingLogger(path string, maxSizeMB int) (*RotatingLogger, error) {
    // TODO: 1. Create/open log file
	file, err := os.OpenFile(path,os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
		
	}
    // TODO: 2. Initialize struct
	maxSize := int64(maxSizeMB) * 1024 * 1024

	logger := &RotatingLogger{
		file: file,
		path: path,
		maxSize: maxSize,
	}
	return logger,nil
}

// Write handles concurrent-safe writing and rotation
func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

	fileInfo, err := rl.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.Size()+ int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("failed to rotate log file: %w", err)
		}
	}
	n, err = rl.file.Write(p)
	if err != nil {
		return n, fmt.Errorf("failed to write to log file: %w", err)
	}

	return
}

// rotate performs atomic log rotation
func (rl *RotatingLogger) rotate() error {
    // TODO: 6. Close current file
	if err := rl.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}
    // TODO: 7. Create timestamped backup
	timestamp := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS
	backupPath := fmt.Sprintf("%s.%s", rl.path, timestamp)
	if err := os.Rename(rl.path, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file to %s: %w", backupPath, err)
	}
    // TODO: 8. Create new log file
	newFile, err := os.OpenFile(rl.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file %s: %w", rl.path, err)
	}

	// Update the logger's file handle to the new file.
	rl.file = newFile

	fmt.Printf("Log file rotated: %s -> %s. New log file created: %s\n", rl.path, backupPath, rl.path)
	rl.logRotation(backupPath)

	return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}
func (rl *RotatingLogger) logRotation(backupPath string) {
    logLine := fmt.Sprintf("[ROTATION] %s rotated to %s\n", 
        rl.path, 
        backupPath,
    )
    // Write through our own Write() method (with mutex already locked)
    rl.file.Write([]byte(logLine))
}