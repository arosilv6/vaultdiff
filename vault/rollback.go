package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Rollbacker rolls a secret back to a previous version.
type Rollbacker struct {
	client *vaultapi.Client
	mount  string
}

// NewRollbacker creates a new Rollbacker.
func NewRollbacker(client *vaultapi.Client, mount string) *Rollbacker {
	return &Rollbacker{client: client, mount: mount}
}

// RollbackResult holds information about a completed rollback.
type RollbackResult struct {
	Path       string
	FromVersion int
	ToVersion   int
}

// Rollback copies the data from toVersion and writes it as a new version.
func (r *Rollbacker) Rollback(ctx context.Context, path string, toVersion int) (*RollbackResult, error) {
	if r.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	secretPath := fmt.Sprintf("%s/data/%s", r.mount, path)
	params := map[string]interface{}{"version": toVersion}

	secret, err := r.client.Logical().ReadWithDataWithContext(ctx, secretPath, params)
	if err != nil {
		return nil, fmt.Errorf("read version %d: %w", toVersion, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("version %d not found at %s", toVersion, path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("no data field in version %d", toVersion)
	}

	writeSecret, err := r.client.Logical().WriteWithContext(ctx, secretPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("write rollback: %w", err)
	}

	newVersion := 0
	if writeSecret != nil && writeSecret.Data != nil {
		if v, ok := writeSecret.Data["version"]; ok {
			if vi, ok := v.(float64); ok {
				newVersion = int(vi)
			}
		}
	}

	return &RollbackResult{Path: path, FromVersion: newVersion, ToVersion: toVersion}, nil
}
