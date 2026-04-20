package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Cloner copies all versions of a secret from one path to another.
type Cloner struct {
	client *api.Client
	mount  string
}

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourcePath string
	DestPath   string
	Versions   []int
}

// NewCloner creates a new Cloner instance.
func NewCloner(client *api.Client, mount string) *Cloner {
	return &Cloner{client: client, mount: mount}
}

// Clone reads all available versions from srcPath and writes them
// sequentially to dstPath, preserving order.
func (c *Cloner) Clone(srcPath, dstPath string) (*CloneResult, error) {
	if c.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	meta, err := c.client.Logical().Read(
		fmt.Sprintf("%s/metadata/%s", c.mount, srcPath),
	)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s: %w", srcPath, err)
	}
	if meta == nil || meta.Data == nil {
		return nil, fmt.Errorf("no metadata found at %s", srcPath)
	}

	versionsRaw, ok := meta.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata for %s", srcPath)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format for %s", srcPath)
	}

	result := &CloneResult{SourcePath: srcPath, DestPath: dstPath}

	for vStr := range versionsMap {
		var vNum int
		fmt.Sscanf(vStr, "%d", &vNum)

		data, err := c.client.Logical().ReadWithData(
			fmt.Sprintf("%s/data/%s", c.mount, srcPath),
			map[string][]string{"version": {vStr}},
		)
		if err != nil || data == nil || data.Data == nil {
			continue
		}

		_, err = c.client.Logical().Write(
			fmt.Sprintf("%s/data/%s", c.mount, dstPath),
			map[string]interface{}{"data": data.Data["data"]},
		)
		if err != nil {
			return nil, fmt.Errorf("writing version %d to %s: %w", vNum, dstPath, err)
		}
		result.Versions = append(result.Versions, vNum)
	}

	return result, nil
}
