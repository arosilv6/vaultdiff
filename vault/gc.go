package vault

import (
	"context"
	"fmt"

	hashivault "github.com/hashicorp/vault/api"
)

// GCResult holds the result of a garbage collection run.
type GCResult struct {
	Path           string
	VersionsPurged []int
	Count          int
}

// GarbageCollector removes old or destroyed secret versions beyond a retention limit.
type GarbageCollector struct {
	client *hashivault.Client
	mount  string
}

// NewGarbageCollector creates a new GarbageCollector.
func NewGarbageCollector(client *hashivault.Client, mount string) *GarbageCollector {
	return &GarbageCollector{client: client, mount: mount}
}

// Collect deletes permanently destroyed versions at path, keeping at most maxVersions.
// Versions beyond maxVersions that are also destroyed are purged via the destroy endpoint.
func (gc *GarbageCollector) Collect(ctx context.Context, path string, maxVersions int) (*GCResult, error) {
	if gc.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if maxVersions < 1 {
		return nil, fmt.Errorf("maxVersions must be at least 1")
	}

	meta, err := gc.client.Logical().ReadWithContext(ctx,
		fmt.Sprintf("%s/metadata/%s", gc.mount, path))
	if err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}
	if meta == nil || meta.Data == nil {
		return &GCResult{Path: path}, nil
	}

	versionsRaw, ok := meta.Data["versions"].(map[string]interface{})
	if !ok {
		return &GCResult{Path: path}, nil
	}

	var toPurge []int
	for vStr, vMeta := range versionsRaw {
		vMap, ok := vMeta.(map[string]interface{})
		if !ok {
			continue
		}
		destroyed, _ := vMap["destroyed"].(bool)
		if !destroyed {
			continue
		}
		var vNum int
		fmt.Sscanf(vStr, "%d", &vNum)
		if vNum > 0 && vNum <= maxVersions {
			toPurge = append(toPurge, vNum)
		}
	}

	if len(toPurge) == 0 {
		return &GCResult{Path: path}, nil
	}

	versionsPayload := make([]interface{}, len(toPurge))
	for i, v := range toPurge {
		versionsPayload[i] = v
	}
	_, err = gc.client.Logical().WriteWithContext(ctx,
		fmt.Sprintf("%s/destroy/%s", gc.mount, path),
		map[string]interface{}{"versions": versionsPayload})
	if err != nil {
		return nil, fmt.Errorf("destroy versions: %w", err)
	}

	return &GCResult{Path: path, VersionsPurged: toPurge, Count: len(toPurge)}, nil
}
