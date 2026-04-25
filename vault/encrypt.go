package vault

import (
	"context"
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Encrypter wraps a Vault client to encrypt and decrypt data via the Transit engine.
type Encrypter struct {
	client  *vaultapi.Client
	mount   string
	keyName string
}

// EncryptResult holds the ciphertext returned by Vault Transit.
type EncryptResult struct {
	KeyName    string
	Ciphertext string
}

// DecryptResult holds the plaintext returned by Vault Transit.
type DecryptResult struct {
	KeyName   string
	Plaintext string
}

// NewEncrypter creates an Encrypter targeting the given Transit mount and key.
func NewEncrypter(client *vaultapi.Client, mount, keyName string) *Encrypter {
	return &Encrypter{client: client, mount: mount, keyName: keyName}
}

// Encrypt base64-encodes plaintext and sends it to the Transit encrypt endpoint.
func (e *Encrypter) Encrypt(ctx context.Context, base64Plaintext string) (*EncryptResult, error) {
	if e.client == nil {
		return nil, errors.New("encrypt: nil vault client")
	}
	path := fmt.Sprintf("%s/encrypt/%s", e.mount, e.keyName)
	secret, err := e.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"plaintext": base64Plaintext,
	})
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("encrypt: empty response from vault")
	}
	ct, _ := secret.Data["ciphertext"].(string)
	return &EncryptResult{KeyName: e.keyName, Ciphertext: ct}, nil
}

// Decrypt sends ciphertext to the Transit decrypt endpoint and returns plaintext.
func (e *Encrypter) Decrypt(ctx context.Context, ciphertext string) (*DecryptResult, error) {
	if e.client == nil {
		return nil, errors.New("decrypt: nil vault client")
	}
	path := fmt.Sprintf("%s/decrypt/%s", e.mount, e.keyName)
	secret, err := e.client.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"ciphertext": ciphertext,
	})
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("decrypt: empty response from vault")
	}
	pt, _ := secret.Data["plaintext"].(string)
	return &DecryptResult{KeyName: e.keyName, Plaintext: pt}, nil
}
