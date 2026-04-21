package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Patcher applies partial updates to a secret version without overwriting unchanged keys.
type Patcher struct {
	client *vaultapi.Client
	mount  string
}

// PatchResult holds the outcome of a patch operation.
type PatchResult struct {
	Path        string
	Version     int
	KeysPatched []string
}

// NewPatcher creates a new Patcher for the given mount.
func NewPatcher(client *vaultapi.Client, mount string) *Patcher {
	return &Patcher{client: client, mount: mount}
}

// Patch reads the latest version of path, merges updates into it, and writes a new version.
// Only keys present in updates are changed; all other keys are preserved.
func (p *Patcher) Patch(ctx context.Context, path string, updates map[string]string) (*PatchResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates provided")
	}

	dataPath := fmt.Sprintf("%s/data/%s", p.mount, path)

	// Read current secret
	secret, err := p.client.Logical().ReadWithContext(ctx, dataPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	current := make(map[string]interface{})
	if secret != nil && secret.Data != nil {
		if data, ok := secret.Data["data"].(map[string]interface{}); ok {
			for k, v := range data {
				current[k] = v
			}
		}
	}

	patched := make([]string, 0, len(updates))
	for k, v := range updates {
		current[k] = v
		patched = append(patched, k)
	}

	writeSecret, err := p.client.Logical().WriteWithContext(ctx, dataPath, map[string]interface{}{
		"data": current,
	})
	if err != nil {
		return nil, fmt.Errorf("write %s: %w", path, err)
	}

	version := 0
	if writeSecret != nil && writeSecret.Data != nil {
		if v, ok := writeSecret.Data["version"].(float64); ok {
			version = int(v)
		}
	}

	return &PatchResult{
		Path:        path,
		Version:     version,
		KeysPatched: patched,
	}, nil
}
