package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Annotation holds a user-defined note attached to a specific secret version.
type Annotation struct {
	Path    string    `json:"path"`
	Version int       `json:"version"`
	Note    string    `json:"note"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
}

// Annotator reads and writes annotations stored as Vault KV metadata
// custom_metadata fields, keyed by version number.
type Annotator struct {
	client *vaultapi.Client
	mount  string
}

// NewAnnotator creates an Annotator targeting the given KV v2 mount.
func NewAnnotator(client *vaultapi.Client, mount string) *Annotator {
	return &Annotator{
		client: client,
		mount:  mount,
	}
}

// SetAnnotation writes a note for the given path and version into the secret's
// custom_metadata. The key format is "note_vN" where N is the version number.
func (a *Annotator) SetAnnotation(ctx context.Context, path string, version int, note, author string) (*Annotation, error) {
	if a.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	if note == "" {
		return nil, fmt.Errorf("annotation note must not be empty")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", a.mount, path)
	noteKey := fmt.Sprintf("note_v%d", version)
	authorKey := fmt.Sprintf("author_v%d", version)

	// Read existing custom_metadata so we don't overwrite other keys.
	existing := map[string]interface{}{}
	secret, err := a.client.Logical().ReadWithContext(ctx, metaPath)
	if err == nil && secret != nil {
		if cm, ok := secret.Data["custom_metadata"].(map[string]interface{}); ok {
			for k, v := range cm {
				existing[k] = v
			}
		}
	}

	existing[noteKey] = note
	existing[authorKey] = author

	_, err = a.client.Logical().WriteWithContext(ctx, metaPath, map[string]interface{}{
		"custom_metadata": existing,
	})
	if err != nil {
		return nil, fmt.Errorf("writing annotation metadata: %w", err)
	}

	return &Annotation{
		Path:    path,
		Version: version,
		Note:    note,
		Author:  author,
		Created: time.Now().UTC(),
	}, nil
}

// GetAnnotation retrieves the annotation for a specific path and version.
// Returns nil without error when no annotation exists for that version.
func (a *Annotator) GetAnnotation(ctx context.Context, path string, version int) (*Annotation, error) {
	if a.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	metaPath := fmt.Sprintf("%s/metadata/%s", a.mount, path)
	secret, err := a.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata for %s: %w", path, err)
	}
	if secret == nil {
		return nil, nil
	}

	cm, ok := secret.Data["custom_metadata"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	noteKey := fmt.Sprintf("note_v%d", version)
	authorKey := fmt.Sprintf("author_v%d", version)

	note, _ := cm[noteKey].(string)
	if note == "" {
		return nil, nil
	}
	author, _ := cm[authorKey].(string)

	return &Annotation{
		Path:    path,
		Version: version,
		Note:    note,
		Author:  author,
	}, nil
}
