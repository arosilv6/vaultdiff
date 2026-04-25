package vault

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// HistoryEntry represents a single version entry in a secret's history.
type HistoryEntry struct {
	Version     int
	CreatedTime time.Time
	Deleted     bool
	Destroyed   bool
}

// Historian retrieves the full version history of a secret path.
type Historian struct {
	client *api.Client
	mount  string
}

// NewHistorian creates a new Historian for the given mount and client.
func NewHistorian(client *api.Client, mount string) *Historian {
	return &Historian{client: client, mount: mount}
}

// History returns all version entries for the given secret path.
func (h *Historian) History(path string) ([]HistoryEntry, error) {
	if h.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", h.mount, path)
	secret, err := h.client.Logical().Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions field in metadata for %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %q", path)
	}

	entries := make([]HistoryEntry, 0, len(versionsMap))
	for _, v := range versionsMap {
		vm, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		entry := HistoryEntry{}
		if ct, ok := vm["created_time"].(string); ok {
			entry.CreatedTime, _ = time.Parse(time.RFC3339Nano, ct)
		}
		if del, ok := vm["deletion_time"].(string); ok && del != "" {
			entry.Deleted = true
		}
		if dest, ok := vm["destroyed"].(bool); ok {
			entry.Destroyed = dest
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
