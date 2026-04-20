package vault

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// Exporter handles exporting Vault secret versions to a file.
type Exporter struct {
	client *api.Client
	mount  string
}

// ExportEntry represents a single exported secret version.
type ExportEntry struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Data    map[string]string `json:"data"`
}

// NewExporter creates a new Exporter instance.
func NewExporter(client *api.Client, mount string) *Exporter {
	return &Exporter{client: client, mount: mount}
}

// ExportVersion reads a specific version of a secret and returns an ExportEntry.
func (e *Exporter) ExportVersion(path string, version int) (*ExportEntry, error) {
	if e.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	versionedPath := fmt.Sprintf("%s/data/%s", e.mount, path)
	secret, err := e.client.Logical().ReadWithData(versionedPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, fmt.Errorf("reading secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at %s version %d", path, version)
	}

	raw, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}

	data := make(map[string]string, len(raw))
	for k, v := range raw {
		data[k] = fmt.Sprintf("%v", v)
	}

	return &ExportEntry{Path: path, Version: version, Data: data}, nil
}

// WriteToFile serialises entries as newline-delimited JSON into the given file path.
func (e *Exporter) WriteToFile(entries []*ExportEntry, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			return fmt.Errorf("encoding entry for %s: %w", entry.Path, err)
		}
	}
	return nil
}
