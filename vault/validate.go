package vault

import (
	"fmt"
	"strings"

	hashivault "github.com/hashicorp/vault/api"
)

// ValidationResult holds the outcome of a secret validation check.
type ValidationResult struct {
	Path    string
	Version int
	Missing []string
	Valid   bool
}

// Validator checks secrets against a set of required keys.
type Validator struct {
	client       *hashivault.Client
	mount        string
	requiredKeys []string
}

// NewValidator creates a Validator for the given mount and required keys.
func NewValidator(client *hashivault.Client, mount string, requiredKeys []string) *Validator {
	return &Validator{
		client:       client,
		mount:        mount,
		requiredKeys: requiredKeys,
	}
}

// Validate reads the secret at path/version and checks for required keys.
// If version is 0, the latest version is read.
func (v *Validator) Validate(path string, version int) (*ValidationResult, error) {
	if v.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	versionParam := ""
	if version > 0 {
		versionParam = fmt.Sprintf("?version=%d", version)
	}

	secretPath := fmt.Sprintf("%s/data/%s%s", strings.TrimRight(v.mount, "/"), strings.TrimLeft(path, "/"), versionParam)
	secret, err := v.client.Logical().Read(secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %s: %w", path, err)
	}

	result := &ValidationResult{Path: path, Version: version}

	if secret == nil || secret.Data == nil {
		result.Missing = v.requiredKeys
		return result, nil
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		result.Missing = v.requiredKeys
		return result, nil
	}

	for _, key := range v.requiredKeys {
		if _, exists := data[key]; !exists {
			result.Missing = append(result.Missing, key)
		}
	}

	result.Valid = len(result.Missing) == 0
	return result, nil
}
