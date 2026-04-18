package audit

import (
	"fmt"
	"time"
)

// WatchRecord represents an audit log entry for a watch event.
type WatchRecord struct {
	Timestamp  time.Time `json:"timestamp"`
	Path       string    `json:"path"`
	OldVersion int       `json:"old_version"`
	NewVersion int       `json:"new_version"`
	Error      string    `json:"error,omitempty"`
}

// RecordWatchEvent writes a watch event to the audit logger.
func (l *Logger) RecordWatchEvent(path string, oldV, newV int, watchErr error) error {
	rec := WatchRecord{
		Timestamp:  time.Now().UTC(),
		Path:       path,
		OldVersion: oldV,
		NewVersion: newV,
	}
	if watchErr != nil {
		rec.Error = watchErr.Error()
	}
	return l.Record(fmt.Sprintf("watch:%s", path), rec)
}
