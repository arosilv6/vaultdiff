package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// Importer reads a JSON export file and writes secrets into Vault.
type Importer struct {
	client *api.Client
	mount  string
}

// ImportEntry represents a single secret entry from an export file.
type ImportEntry struct {
	Path string                 `json:"path"`
	Data map[string]interface{} `json:"data"`
}

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Imported int
	Skipped  int
	Errors   []string
}

// NewImporter creates a new Importer for the given mount.
func NewImporter(client *api.Client, mount string) *Importer {
	return &Importer{client: client, mount: mount}
}

// ImportFile reads entries from a JSON file and writes them to Vault.
func (i *Importer) ImportFile(ctx context.Context, filePath string, dryRun bool) (*ImportResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var entries []ImportEntry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	result := &ImportResult{}
	for _, entry := range entries {
		if i.client == nil {
			result.Errors = append(result.Errors, fmt.Sprintf("nil client for path %s", entry.Path))
			result.Skipped++
			continue
		}
		if dryRun {
			result.Imported++
			continue
		}
		secretPath := fmt.Sprintf("%s/data/%s", i.mount, entry.Path)
		_, err := i.client.Logical().WriteWithContext(ctx, secretPath, map[string]interface{}{
			"data": entry.Data,
		})
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", entry.Path, err))
			result.Skipped++
			continue
		}
		result.Imported++
	}
	return result, nil
}
