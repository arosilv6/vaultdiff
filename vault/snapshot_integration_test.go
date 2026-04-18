package vault

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestSnapshot_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("VAULT_ADDR or VAULT_TOKEN not set")
	}

	cfg := api.DefaultConfig()
	cfg.Address = addr
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}
	client.SetToken(token)

	s := NewSnapshotter(client, "secret")
	entries, err := s.Snapshot("")
	if err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}
	t.Logf("snapshot captured %d entries", len(entries))
}
