package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Comparer reads two versions of a KV secret for diffing.
type Comparer struct {
	client *api.Client
	mount  string
}

// CompareResult holds the raw data for two versions of a secret.
type CompareResult struct {
	Path    string
	VersionA int
	VersionB int
	DataA   map[string]interface{}
	DataB   map[string]interface{}
}

// NewComparer creates a new Comparer.
func NewComparer(client *api.Client, mount string) *Comparer {
	return &Comparer{client: client, mount: mount}
}

// FetchVersions retrieves data for two versions of a secret path.
func (c *Comparer) FetchVersions(path string, versionA, versionB int) (*CompareResult, error) {
	if c.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	dataA, err := c.readVersion(path, versionA)
	if err != nil {
		return nil, fmt.Errorf("reading version %d: %w", versionA, err)
	}

	dataB, err := c.readVersion(path, versionB)
	if err != nil {
		return nil, fmt.Errorf("reading version %d: %w", versionB, err)
	}

	return &CompareResult{
		Path:     path,
		VersionA: versionA,
		VersionB: versionB,
		DataA:    dataA,
		DataB:    dataB,
	}, nil
}

func (c *Comparer) readVersion(path string, version int) (map[string]interface{}, error) {
	kvPath := fmt.Sprintf("%s/data/%s", c.mount, path)
	secret, err := c.client.Logical().ReadWithData(kvPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return map[string]interface{}{}, nil
	}
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}
	return data, nil
}
