package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Tagger manages human-readable tags stored as custom metadata on KV v2 secrets.
type Tagger struct {
	client *api.Client
	mount  string
}

// NewTagger creates a new Tagger for the given mount.
func NewTagger(client *api.Client, mount string) *Tagger {
	return &Tagger{client: client, mount: mount}
}

// SetTag writes a tag (key=value) into the secret's custom_metadata.
func (t *Tagger) SetTag(ctx context.Context, path, key, value string) error {
	if t.client == nil {
		return fmt.Errorf("vault client is nil")
	}
	metaPath := fmt.Sprintf("%s/metadata/%s", t.mount, path)
	body := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			key: value,
		},
	}
	_, err := t.client.Logical().WriteWithContext(ctx, metaPath, body)
	if err != nil {
		return fmt.Errorf("set tag: %w", err)
	}
	return nil
}

// GetTags returns the custom_metadata map for the given secret path.
func (t *Tagger) GetTags(ctx context.Context, path string) (map[string]string, error) {
	if t.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}
	metaPath := fmt.Sprintf("%s/metadata/%s", t.mount, path)
	secret, err := t.client.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("get tags: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return map[string]string{}, nil
	}
	raw, ok := secret.Data["custom_metadata"]
	if !ok {
		return map[string]string{}, nil
	}
	cm, ok := raw.(map[string]interface{})
	if !ok {
		return map[string]string{}, nil
	}
	tags := make(map[string]string, len(cm))
	for k, v := range cm {
		tags[k] = fmt.Sprintf("%v", v)
	}
	return tags, nil
}
