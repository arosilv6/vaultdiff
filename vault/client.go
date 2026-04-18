package vault

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	api *vaultapi.Client
}

// NewClient creates a new Vault client using environment variables
// (VAULT_ADDR, VAULT_TOKEN, etc.).
func NewClient() (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	if err := cfg.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("reading vault environment: %w", err)
	}

	c, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	return &Client{api: c}, nil
}

// GetSecretVersion retrieves a specific version of a KV v2 secret.
func (c *Client) GetSecretVersion(mount, path string, version int) (map[string]string, error) {
	kvPath := fmt.Sprintf("%s/data/%s", mount, path)

	secret, err := c.api.Logical().ReadWithData(kvPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, fmt.Errorf("reading secret %s@v%d: %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %s version %d not found", path, version)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for %s", path)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
