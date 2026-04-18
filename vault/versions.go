package vault

import (
	"fmt"
	"sort"
)

// VersionMeta holds metadata about a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KV v2 secret.
func (c *Client) ListVersions(mountPath, secretPath string) ([]VersionMeta, error) {
	path := fmt.Sprintf("%s/metadata/%s", mountPath, secretPath)

	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading metadata at %s: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found at %s", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata response")
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format")
	}

	var versions []VersionMeta
	for _, v := range versionsMap {
		data, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{}
		if ct, ok := data["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := data["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if destroyed, ok := data["destroyed"].(bool); ok {
			meta.Destroyed = destroyed
		}
		versions = append(versions, meta)
	}

	// Assign version numbers by sorted index (Vault returns them as string keys)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedTime < versions[j].CreatedTime
	})
	for i := range versions {
		versions[i].Version = i + 1
	}

	return versions, nil
}
