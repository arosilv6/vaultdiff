package vault

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestCompare_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Skip("VAULT_ADDR or VAULT_TOKEN not set; skipping integration test")
	}

	cfg := api.DefaultConfig()
	cfg.Address = addr
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("creating vault client: %v", err)
	}
	client.SetToken(token)

	comparer := NewComparer(client, "secret")
	if comparer == nil {
		t.Fatal("expected non-nil Comparer")
	}

	result, err := comparer.FetchVersions("vaultdiff/integration", 1, 2)
	if err != nil {
		t.Logf("FetchVersions returned error (expected if path absent): %v", err)
		return
	}
	if result.Path != "vaultdiff/integration" {
		t.Errorf("unexpected path: %s", result.Path)
	}
	if result.VersionA != 1 || result.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", result.VersionA, result.VersionB)
	}
	t.Logf("DataA: %v", result.DataA)
	t.Logf("DataB: %v", result.DataB)
}
