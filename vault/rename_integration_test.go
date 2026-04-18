//go:build integration
// +build integration

package vault_test

import (
	"os"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/user/vaultdiff/vault"
)

func TestRename_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("VAULT_ADDR or VAULT_TOKEN not set")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	client.SetToken(token)

	// Write initial secret
	_, err = client.Logical().Write("secret/data/rename-src", map[string]interface{}{
		"data": map[string]interface{}{"key": "value"},
	})
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	renamer := vault.NewRenamer(client, "secret")
	if err := renamer.Rename("rename-src", "rename-dst"); err != nil {
		t.Fatalf("rename: %v", err)
	}

	// Verify destination exists
	s, err := client.Logical().Read("secret/data/rename-dst")
	if err != nil || s == nil {
		t.Fatal("destination secret not found after rename")
	}

	// Cleanup
	client.Logical().Delete("secret/metadata/rename-dst")
}
