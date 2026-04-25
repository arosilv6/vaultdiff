package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// ExpiryInfo holds TTL metadata for a secret path.
type ExpiryInfo struct {
	Path      string
	Version   int
	ExpiresAt time.Time
	TTL       time.Duration
	Expired   bool
}

// Expirer checks secret lease/TTL expiry information.
type Expirer struct {
	client *api.Client
	mount  string
}

// NewExpirer creates a new Expirer.
func NewExpirer(client *api.Client, mount string) *Expirer {
	return &Expirer{client: client, mount: mount}
}

// CheckExpiry returns expiry info for the given secret path and version.
// It reads KV v2 metadata from Vault and parses the deletion_time field
// for the specified version. If no deletion_time is set, ExpiresAt will
// be the zero value and TTL will be zero (indicating no expiry configured).
func (e *Expirer) CheckExpiry(ctx context.Context, path string, version int) (*ExpiryInfo, error) {
	if e.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", e.mount, path)
	secret, err := e.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for %s", path)
	}

	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format")
	}

	key := fmt.Sprintf("%d", version)
	vMeta, ok := versions[key].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("version %d not found", version)
	}

	info := &ExpiryInfo{Path: path, Version: version}
	if delTime, ok := vMeta["deletion_time"].(string); ok && delTime != "" {
		parsed, err := time.Parse(time.RFC3339, delTime)
		if err != nil {
			return nil, fmt.Errorf("parsing deletion_time %q: %w", delTime, err)
		}
		info.ExpiresAt = parsed
		info.TTL = time.Until(parsed)
		info.Expired = time.Now().After(parsed)
	}

	return info, nil
}
