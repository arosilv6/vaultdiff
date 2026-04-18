package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/vaultdiff/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Path      string           `json:"path"`
	VersionA  int              `json:"version_a"`
	VersionB  int              `json:"version_b"`
	Changes   []diff.Change    `json:"changes"`
}

// Logger writes audit entries to a file or stdout.
type Logger struct {
	out *os.File
}

// NewLogger creates a Logger. If path is empty, stdout is used.
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{out: os.Stdout}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{out: f}, nil
}

// Record writes an audit entry as a JSON line.
func (l *Logger) Record(path string, versionA, versionB int, changes []diff.Change) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Path:      path,
		VersionA:  versionA,
		VersionB:  versionB,
		Changes:   changes,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(l.out, string(data))
	return err
}

// Close closes the underlying file if it is not stdout.
func (l *Logger) Close() error {
	if l.out != os.Stdout {
		return l.out.Close()
	}
	return nil
}
