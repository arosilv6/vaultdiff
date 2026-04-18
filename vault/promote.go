package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Promoter copies a specific version of a secret to the latest version,
// effectively "promoting" an older version to current.
type Promoter struct {
	client *api.Client
	mount  string
}

// NewPromoter creates a new Promoter for the given mount.
func NewPromoter(client *api.Client, mount string) *Promoter {
	return &Promoter{client: client, mount: mount}
}

// Promote reads the given version of path and writes it as a new version.
func (p *Promoter) Promote(path string, version int) (map[string]interface{}, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	fullPath := fmt.Sprintf("%s/data/%s", p.mount, path)

	// Read the specific version
	secret, err := p.client.Logical().ReadWithData(fullPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, fmt.Errorf("reading version %d of %s: %w", version, path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at %s version %d", path, version)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at %s", path)
	}

	// Write as new version
	result, err := p.client.Logical().Write(fullPath, map[string]interface{}{
		"data": data,
	})
	if err != nil {
		return nil, fmt.Errorf("promoting version %d of %s: %w", version, path, err)
	}
	if result == nil {
		return data, nil
	}
	return data, nil
}
