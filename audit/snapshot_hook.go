package audit

import (
	"time"
)

// SnapshotRecord is written to the audit log when a snapshot is taken.
type SnapshotRecord struct {
	Timestamp  time.Time `json:"timestamp"`
	Action     string    `json:"action"`
	Mount      string    `json:"mount"`
	BasePath   string    `json:"base_path"`
	EntryCount int       `json:"entry_count"`
}

// RecordSnapshot logs a snapshot operation.
func (l *Logger) RecordSnapshot(mount, basePath string, count int) error {
	return l.Record(SnapshotRecord{
		Timestamp:  time.Now().UTC(),
		Action:     "snapshot",
		Mount:      mount,
		BasePath:   basePath,
		EntryCount: count,
	})
}
