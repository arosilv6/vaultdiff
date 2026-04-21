package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// BlameEntry represents who wrote a specific version of a secret.
type BlameEntry struct {
	Version     int
	CreatedTime time.Time
	CreatedBy   string
	Operation   string
	Deleted     bool
}

// Blamer retrieves authorship history for a Vault KV secret path.
type Blamer struct {
	client *vaultapi.Client
	mount  string
}

// NewBlamer creates a new Blamer instance.
func NewBlamer(client *vaultapi.Client, mount string) *Blamer {
	return &Blamer{client: client, mount: mount}
}

// Blame returns a BlameEntry slice for all versions of the given path.
func (b *Blamer) Blame(ctx context.Context, path string) ([]BlameEntry, error) {
	if b.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", b.mount, path)
	secret, err := b.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path: %s", path)
	}

	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for path: %s", path)
	}

	entries := make([]BlameEntry, 0, len(versions))
	for versionKey, raw := range versions {
		vMeta, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		var vNum int
		fmt.Sscanf(versionKey, "%d", &vNum)

		entry := BlameEntry{Version: vNum}

		if ct, ok := vMeta["created_time"].(string); ok {
			entry.CreatedTime, _ = time.Parse(time.RFC3339Nano, ct)
		}
		if by, ok := vMeta["created_by"].(string); ok {
			entry.CreatedBy = by
		}
		if op, ok := vMeta["operation"].(string); ok {
			entry.Operation = op
		}
		if dt, ok := vMeta["deletion_time"].(string); ok && dt != "" {
			entry.Deleted = true
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
