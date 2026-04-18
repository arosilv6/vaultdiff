package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Restorer handles rolling back secrets to a previous version.
type Restorer struct {
	client *api.Client
	mount  string
}

// NewRestorer creates a Restorer for the given mount path.
func NewRestorer(client *api.Client, mount string) *Restorer {
	return &Restorer{client: client, mount: mount}
}

// Rollback sets the current version of the secret at path to the given version.
func (r *Restorer) Rollback(path string, version int) error {
	kvPath := fmt.Sprintf("%s/undelete/%s", r.mount, path)
	body := map[string]interface{}{
		"versions": []int{version},
	}
	_, err := r.client.Logical().Write(kvPath, body)
	if err != nil {
		return fmt.Errorf("undelete version %d at %s: %w", version, path, err)
	}

	kvPath = fmt.Sprintf("%s/metadata/%s", r.mount, path)
	body = map[string]interface{}{
		"current_version": version,
	}
	_, err = r.client.Logical().Write(kvPath, body)
	if err != nil {
		return fmt.Errorf("set current version %d at %s: %w", version, path, err)
	}
	return nil
}

// ReadVersion reads a specific version of a secret.
func (r *Restorer) ReadVersion(path string, version int) (map[string]interface{}, error) {
	kvPath := fmt.Sprintf("%s/data/%s", r.mount, path)
	secret, err := r.client.Logical().ReadWithData(kvPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, fmt.Errorf("read version %d at %s: %w", version, path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found for version %d at %s", version, path)
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at %s", path)
	}
	return data, nil
}
