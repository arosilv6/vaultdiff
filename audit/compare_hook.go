package audit

import (
	"time"
)

// CompareEvent records a diff comparison audit event.
type CompareEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Path      string    `json:"path"`
	Mount     string    `json:"mount"`
	VersionA  int       `json:"version_a"`
	VersionB  int       `json:"version_b"`
	Changes   int       `json:"changes"`
}

// RecordCompare logs a compare action to the audit logger.
func RecordCompare(logger *Logger, path, mount string, versionA, versionB, changes int) error {
	event := CompareEvent{
		Timestamp: time.Now().UTC(),
		Action:    "compare",
		Path:      path,
		Mount:     mount,
		VersionA:  versionA,
		VersionB:  versionB,
		Changes:   changes,
	}
	return logger.Record(event)
}
