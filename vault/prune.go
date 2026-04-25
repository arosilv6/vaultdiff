package vault

import (
	"context"
	"fmt"

	hashivault "github.com/hashicorp/vault/api"
)

// PruneResult holds the outcome of a prune operation.
type PruneResult struct {
	Path            string
	VersionsPruned  []int
	VersionsKept    int
	DryRun          bool
}

// Pruner removes old secret versions beyond a configurable keep count.
type Pruner struct {
	client    *hashivault.Client
	mountPath string
}

// NewPruner creates a new Pruner for the given mount path.
func NewPruner(client *hashivault.Client, mountPath string) *Pruner {
	return &Pruner{client: client, mountPath: mountPath}
}

// Prune deletes all versions of the secret at path, keeping the newest `keep` versions.
// If dryRun is true, it reports what would be deleted without making changes.
func (p *Pruner) Prune(ctx context.Context, secretPath string, keep int, dryRun bool) (*PruneResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if keep < 1 {
		return nil, fmt.Errorf("keep must be at least 1")
	}

	meta, err := p.client.Logical().ReadWithContext(ctx,
		fmt.Sprintf("%s/metadata/%s", p.mountPath, secretPath))
	if err != nil {
		return nil, fmt.Errorf("reading metadata: %w", err)
	}
	if meta == nil || meta.Data == nil {
		return nil, fmt.Errorf("no metadata found for %s", secretPath)
	}

	versionsRaw, ok := meta.Data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected versions format")
	}

	allVersions := make([]int, 0, len(versionsRaw))
	for k := range versionsRaw {
		var v int
		fmt.Sscanf(k, "%d", &v)
		allVersions = append(allVersions, v)
	}

	// Sort ascending
	for i := 0; i < len(allVersions); i++ {
		for j := i + 1; j < len(allVersions); j++ {
			if allVersions[j] < allVersions[i] {
				allVersions[i], allVersions[j] = allVersions[j], allVersions[i]
			}
		}
	}

	var toPrune []int
	if len(allVersions) > keep {
		toPrune = allVersions[:len(allVersions)-keep]
	}

	result := &PruneResult{
		Path:           secretPath,
		VersionsPruned: toPrune,
		VersionsKept:   len(allVersions) - len(toPrune),
		DryRun:         dryRun,
	}

	if dryRun || len(toPrune) == 0 {
		return result, nil
	}

	versionsPayload := make([]interface{}, len(toPrune))
	for i, v := range toPrune {
		versionsPayload[i] = v
	}
	_, err = p.client.Logical().WriteWithContext(ctx,
		fmt.Sprintf("%s/destroy/%s", p.mountPath, secretPath),
		map[string]interface{}{"versions": versionsPayload})
	if err != nil {
		return nil, fmt.Errorf("destroying versions: %w", err)
	}

	return result, nil
}
