package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// PolicyChecker checks whether the current Vault token has read/write access
// to a given secret path by inspecting capabilities.
type PolicyChecker struct {
	client *vaultapi.Client
}

// CapabilityResult holds the capabilities returned for a path.
type CapabilityResult struct {
	Path         string
	Capabilities []string
}

// HasCapability returns true if the given capability (e.g. "read") is present.
func (r CapabilityResult) HasCapability(cap string) bool {
	for _, c := range r.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}

// NewPolicyChecker creates a PolicyChecker using the provided Vault client.
func NewPolicyChecker(client *vaultapi.Client) *PolicyChecker {
	return &PolicyChecker{client: client}
}

// CheckPath queries Vault for the token's capabilities on the given path.
func (p *PolicyChecker) CheckPath(ctx context.Context, path string) (*CapabilityResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	body := map[string]interface{}{
		"path": path,
	}

	secret, err := p.client.Logical().WriteWithContext(ctx, "sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("capabilities check failed: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return &CapabilityResult{Path: path, Capabilities: []string{}}, nil
	}

	raw, ok := secret.Data[path]
	if !ok {
		return &CapabilityResult{Path: path, Capabilities: []string{}}, nil
	}

	caps := []string{}
	if items, ok := raw.([]interface{}); ok {
		for _, item := range items {
			if s, ok := item.(string); ok {
				caps = append(caps, s)
			}
		}
	}

	return &CapabilityResult{Path: path, Capabilities: caps}, nil
}
