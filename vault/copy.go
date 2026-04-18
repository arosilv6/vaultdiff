package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Copier copies a secret from one path/version to another path.
type Copier struct {
	client *vaultapi.Client
	mount  string
}

// NewCopier creates a new Copier.
func NewCopier(client *vaultapi.Client, mount string) *Copier {
	return &Copier{client: client, mount: mount}
}

// CopyVersion reads the given version from srcPath and writes it to dstPath.
func (c *Copier) CopyVersion(ctx context.Context, srcPath string, version int, dstPath string) error {
	readPath := fmt.Sprintf("%s/data/%s?version=%d", c.mount, srcPath, version)
	secret, err := c.client.Logical().ReadWithContext(ctx, readPath)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("no data at %s version %d", srcPath, version)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected data format at %s", srcPath)
	}

	writePath := fmt.Sprintf("%s/data/%s", c.mount, dstPath)
	_, err = c.client.Logical().WriteWithContext(ctx, writePath, map[string]interface{}{
		"data": data,
	})
	if err != nil {
		return fmt.Errorf("writing destination: %w", err)
	}
	return nil
}
