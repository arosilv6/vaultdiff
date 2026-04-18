package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Locker manages advisory locks on Vault KV paths using metadata custom_metadata.
type Locker struct {
	client *api.Client
	mount  string
}

// NewLocker creates a new Locker for the given mount.
func NewLocker(client *api.Client, mount string) *Locker {
	return &Locker{client: client, mount: mount}
}

// Lock sets a "locked" flag in the secret's custom_metadata.
func (l *Locker) Lock(ctx context.Context, path string) error {
	if l.client == nil {
		return fmt.Errorf("vault client is nil")
	}
	metaPath := fmt.Sprintf("%s/metadata/%s", l.mount, path)
	_, err := l.client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"locked": "true",
		},
	})
	if err != nil {
		return fmt.Errorf("lock %s: %w", path, err)
	}
	return nil
}

// Unlock removes the "locked" flag from the secret's custom_metadata.
func (l *Locker) Unlock(ctx context.Context, path string) error {
	if l.client == nil {
		return fmt.Errorf("vault client is nil")
	}
	metaPath := fmt.Sprintf("%s/metadata/%s", l.mount, path)
	_, err := l.client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"locked": "false",
		},
	})
	if err != nil {
		return fmt.Errorf("unlock %s: %w", path, err)
	}
	return nil
}

// IsLocked reads metadata and returns whether the path is locked.
func (l *Locker) IsLocked(ctx context.Context, path string) (bool, error) {
	if l.client == nil {
		return false, fmt.Errorf("vault client is nil")
	}
	metaPath := fmt.Sprintf("%s/metadata/%s", l.mount, path)
	secret, err := l.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return false, fmt.Errorf("read metadata %s: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return false, nil
	}
	cm, ok := secret.Data["custom_metadata"].(map[string]interface{})
	if !ok {
		return false, nil
	}
	return cm["locked"] == "true", nil
}
