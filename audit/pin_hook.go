package audit

import (
	"fmt"
	"time"
)

// PinEvent represents an audit log entry for a pin or unpin action.
type PinEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // "pin" or "unpin"
	Path      string    `json:"path"`
	Version   int       `json:"version,omitempty"`
	Mount     string    `json:"mount"`
}

// RecordPin writes a pin audit event to the provided Logger.
func RecordPin(logger *Logger, mount, path string, version int) error {
	if logger == nil {
		return fmt.Errorf("audit logger is nil")
	}
	return logger.Record(PinEvent{
		Timestamp: time.Now().UTC(),
		Action:    "pin",
		Path:      path,
		Version:   version,
		Mount:     mount,
	})
}

// RecordUnpin writes an unpin audit event to the provided Logger.
func RecordUnpin(logger *Logger, mount, path string) error {
	if logger == nil {
		return fmt.Errorf("audit logger is nil")
	}
	return logger.Record(PinEvent{
		Timestamp: time.Now().UTC(),
		Action:    "unpin",
		Path:      path,
		Mount:     mount,
	})
}
