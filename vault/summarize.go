package vault

import (
	"fmt"
	"sort"

	hashivault "github.com/hashicorp/vault/api"
)

// SummaryEntry holds aggregated metadata for a single secret path.
type SummaryEntry struct {
	Path        string
	VersionCount int
	LatestVersion int
	DeletedCount  int
	DestroyedCount int
}

// Summarizer aggregates version metadata for one or more secret paths.
type Summarizer struct {
	client *hashivault.Client
	mount  string
}

// NewSummarizer constructs a Summarizer for the given mount.
func NewSummarizer(client *hashivault.Client, mount string) *Summarizer {
	return &Summarizer{client: client, mount: mount}
}

// Summarize returns a SummaryEntry for the given secret path.
func (s *Summarizer) Summarize(path string) (*SummaryEntry, error) {
	if s.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	logical := s.client.Logical()
	fullPath := fmt.Sprintf("%s/metadata/%s", s.mount, path)

	secret, err := logical.Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path: %s", path)
	}

	entry := &SummaryEntry{Path: path}

	if v, ok := secret.Data["current_version"]; ok {
		if n, ok := v.(float64); ok {
			entry.LatestVersion = int(n)
		}
	}

	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return entry, nil
	}

	entry.VersionCount = len(versions)

	keys := make([]string, 0, len(versions))
	for k := range versions {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		meta, ok := versions[k].(map[string]interface{})
		if !ok {
			continue
		}
		if dt, ok := meta["deletion_time"].(string); ok && dt != "" {
			entry.DeletedCount++
		}
		if destroyed, ok := meta["destroyed"].(bool); ok && destroyed {
			entry.DestroyedCount++
		}
	}

	return entry, nil
}
