package vault

import (
	"context"
	"fmt"

	hashivault "github.com/hashicorp/vault/api"
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Path        string
	BaseVersion int
	OtherVersion int
	MergedKeys  []string
	Conflicts   []string
}

// Merger merges two versions of a secret into a new version.
type Merger struct {
	client *hashivault.Client
	mount  string
}

// NewMerger creates a new Merger.
func NewMerger(client *hashivault.Client, mount string) *Merger {
	return &Merger{client: client, mount: mount}
}

// Merge reads two versions of a secret and writes their combined data as a new version.
// Keys present in both versions are taken from the "other" version (other wins on conflict).
func (m *Merger) Merge(ctx context.Context, path string, baseVersion, otherVersion int) (*MergeResult, error) {
	if m.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	readVersion := func(ver int) (map[string]interface{}, error) {
		p := fmt.Sprintf("%s/data/%s", m.mount, path)
		secret, err := m.client.Logical().ReadWithDataWithContext(ctx, p,
			map[string][]string{"version": {fmt.Sprintf("%d", ver)}})
		if err != nil {
			return nil, fmt.Errorf("reading version %d: %w", ver, err)
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

	baseData, err := readVersion(baseVersion)
	if err != nil {
		return nil, err
	}
	otherData, err := readVersion(otherVersion)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]interface{})
	var mergedKeys, conflicts []string

	for k, v := range baseData {
		merged[k] = v
		mergedKeys = append(mergedKeys, k)
	}
	for k, v := range otherData {
		if _, exists := merged[k]; exists {
			conflicts = append(conflicts, k)
		} else {
			mergedKeys = append(mergedKeys, k)
		}
		merged[k] = v
	}

	p := fmt.Sprintf("%s/data/%s", m.mount, path)
	_, err = m.client.Logical().WriteWithContext(ctx, p, map[string]interface{}{"data": merged})
	if err != nil {
		return nil, fmt.Errorf("writing merged secret: %w", err)
	}

	return &MergeResult{
		Path:         path,
		BaseVersion:  baseVersion,
		OtherVersion: otherVersion,
		MergedKeys:   mergedKeys,
		Conflicts:    conflicts,
	}, nil
}
