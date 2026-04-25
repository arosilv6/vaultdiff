package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Templater renders secret values into a template string.
type Templater struct {
	client *vaultapi.Client
	mount  string
}

// TemplateResult holds the rendered output and metadata.
type TemplateResult struct {
	Rendered string
	SecretsUsed []string
}

// NewTemplater creates a new Templater.
func NewTemplater(client *vaultapi.Client, mount string) *Templater {
	return &Templater{client: client, mount: mount}
}

// Render replaces {{ secret "path" "key" }} placeholders with actual secret values.
func (t *Templater) Render(ctx context.Context, tmpl string) (*TemplateResult, error) {
	if t.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	result := &TemplateResult{Rendered: tmpl}
	seen := map[string]bool{}

	for {
		start := strings.Index(result.Rendered, "{{ secret \"")
		if start == -1 {
			break
		}
		end := strings.Index(result.Rendered[start:], "}}")
		if end == -1 {
			return nil, fmt.Errorf("unclosed template directive")
		}
		end += start + 2

		directive := result.Rendered[start:end]
		path, key, err := parseDirective(directive)
		if err != nil {
			return nil, err
		}

		secret, err := t.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/data/%s", t.mount, path))
		if err != nil {
			return nil, fmt.Errorf("reading secret %q: %w", path, err)
		}
		if secret == nil || secret.Data == nil {
			return nil, fmt.Errorf("secret %q not found", path)
		}
		data, _ := secret.Data["data"].(map[string]interface{})
		val, ok := data[key]
		if !ok {
			return nil, fmt.Errorf("key %q not found in secret %q", key, path)
		}

		ref := path + ":" + key
		if !seen[ref] {
			result.SecretsUsed = append(result.SecretsUsed, ref)
			seen[ref] = true
		}
		result.Rendered = result.Rendered[:start] + fmt.Sprintf("%v", val) + result.Rendered[end:]
	}

	return result, nil
}

func parseDirective(d string) (path, key string, err error) {
	d = strings.TrimPrefix(d, "{{ secret \"")
	d = strings.TrimSuffix(d, " }}")
	parts := strings.SplitN(d, "\" \"", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid directive format: %q", d)
	}
	return parts[0], strings.TrimSuffix(parts[1], "\""), nil
}
