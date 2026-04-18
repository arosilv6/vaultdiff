package vault_test

import (
	"os"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultdiff/vault"
)

func TestPromote_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("VAULT_ADDR or VAULT_TOKEN not set; skipping integration test")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}
	client.SetToken(token)

	const mount = "secret"
	const path = "vaultdiff-promote-test"

	// Write initial version
	_, err = client.Logical().Write(mount+"/data/"+path, map[string]interface{}{
		"data": map[string]interface{}{"key": "original"},
	})
	if err != nil {
		t.Fatalf("writing v1: %v", err)
	}

	// Write second version
	_, err = client.Logical().Write(mount+"/data/"+path, map[string]interface{}{
		"data": map[string]interface{}{"key": "updated"},
	})
	if err != nil {
		t.Fatalf("writing v2: %v", err)
	}

	promoter := vault.NewPromoter(client, mount)
	data, err := promoter.Promote(path, 1)
	if err != nil {
		t.Fatalf("promote failed: %v", err)
	}

	if data["key"] != "original" {
		t.Errorf("expected promoted data key=original, got %v", data["key"])
	}
}
