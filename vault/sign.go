package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// Signer computes a deterministic HMAC-SHA256 signature over a secret version's
// key/value pairs so callers can detect out-of-band tampering.
type Signer struct {
	client *vaultapi.Client
}

// SignResult holds the path, version, and computed signature.
type SignResult struct {
	Path      string
	Version   int
	Signature string
}

// NewSigner returns a Signer backed by the provided Vault client.
func NewSigner(client *vaultapi.Client) *Signer {
	return &Signer{client: client}
}

// Sign reads the given KV v2 path at the specified version and returns an
// HMAC-SHA256 signature derived from the secret data and the provided key.
// Version <= 0 reads the latest version.
func (s *Signer) Sign(mount, path string, version int, hmacKey string) (*SignResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	kvPath := fmt.Sprintf("%s/data/%s", mount, path)
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := s.client.Logical().ReadWithData(kvPath, params)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", kvPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data at %s", kvPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at %s", kvPath)
	}

	sig := computeSignature(data, hmacKey)

	actualVersion := version
	if meta, ok := secret.Data["metadata"].(map[string]interface{}); ok {
		if v, ok := meta["version"].(json.Number); ok {
			if n, err := v.Int64(); err == nil {
				actualVersion = int(n)
			}
		}
	}

	return &SignResult{
		Path:      path,
		Version:   actualVersion,
		Signature: sig,
	}, nil
}

// computeSignature builds a canonical string from sorted key=value pairs and
// returns its HMAC-SHA256 hex digest.
func computeSignature(data map[string]interface{}, key string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, data[k]))
	}
	canonical := strings.Join(parts, "\n")

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}
