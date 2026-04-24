package vault

import (
	"context"
	"fmt"

	hashivault "github.com/hashicorp/vault/api"
)

// PinResult holds the result of a pin or unpin operation.
type PinResult struct {
	Path    string
	Version int
	Pinned  bool
}

// Pinner manages version pinning for Vault KV secrets.
type Pinner struct {
	client *hashivault.Client
	mount  string
}

// NewPinner creates a new Pinner.
func NewPinner(client *hashivault.Client, mount string) *Pinner {
	return &Pinner{client: client, mount: mount}
}

// Pin marks a specific version of a secret as the pinned version
// by writing a custom metadata key "pinned_version".
func (p *Pinner) Pin(ctx context.Context, path string, version int) (*PinResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if version < 1 {
		return nil, fmt.Errorf("version must be >= 1, got %d", version)
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", p.mount, path)
	_, err := p.client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"pinned_version": fmt.Sprintf("%d", version),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("pin write failed: %w", err)
	}
	return &PinResult{Path: path, Version: version, Pinned: true}, nil
}

// Unpin removes the pinned_version custom metadata key.
func (p *Pinner) Unpin(ctx context.Context, path string) (*PinResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", p.mount, path)
	_, err := p.client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"pinned_version": "",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unpin write failed: %w", err)
	}
	return &PinResult{Path: path, Version: 0, Pinned: false}, nil
}
