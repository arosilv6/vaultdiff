package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Snapshotter captures all secrets under a path at a point in time.
type Snapshotter struct {
	client *api.Client
	mount  string
}

// SnapshotEntry holds a single secret path and its data.
type SnapshotEntry struct {
	Path    string
	Version int
	Data    map[string]interface{}
}

// NewSnapshotter creates a Snapshotter for the given mount.
func NewSnapshotter(client *api.Client, mount string) *Snapshotter {
	return &Snapshotter{client: client, mount: mount}
}

// Snapshot walks all secrets under basePath and returns their current data.
func (s *Snapshotter) Snapshot(basePath string) ([]SnapshotEntry, error) {
	if s.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	searcher := NewSearcher(s.client, s.mount)
	results, err := searcher.ListRecursive(basePath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	var entries []SnapshotEntry
	for _, r := range results {
		secret, err := s.client.Logical().Read(
			fmt.Sprintf("%s/data/%s", s.mount, r.Path),
		)
		if err != nil || secret == nil {
			continue
		}
		data, _ := secret.Data["data"].(map[string]interface{})
		meta, _ := secret.Data["metadata"].(map[string]interface{})
		version := 0
		if v, ok := meta["version"].(float64); ok {
			version = int(v)
		}
		entries = append(entries, SnapshotEntry{
			Path:    r.Path,
			Version: version,
			Data:    data,
		})
	}
	return entries, nil
}
