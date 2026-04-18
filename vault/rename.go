package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Renamer handles moving/renaming secrets within Vault KV v2.
type Renamer struct {
	client *api.Client
	mount  string
}

// NewRenamer creates a new Renamer instance.
func NewRenamer(client *api.Client, mount string) *Renamer {
	return &Renamer{client: client, mount: mount}
}

// Rename copies all data from srcPath to dstPath and then deletes srcPath.
func (r *Renamer) Rename(srcPath, dstPath string) error {
	readPath := fmt.Sprintf("%s/data/%s", r.mount, srcPath)
	secret, err := r.client.Logical().Read(readPath)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("source path %q not found", srcPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected data format at %q", srcPath)
	}

	writePath := fmt.Sprintf("%s/data/%s", r.mount, dstPath)
	_, err = r.client.Logical().Write(writePath, map[string]interface{}{"data": data})
	if err != nil {
		return fmt.Errorf("write destination: %w", err)
	}

	deletePath := fmt.Sprintf("%s/metadata/%s", r.mount, srcPath)
	_, err = r.client.Logical().Delete(deletePath)
	if err != nil {
		return fmt.Errorf("delete source: %w", err)
	}

	return nil
}
