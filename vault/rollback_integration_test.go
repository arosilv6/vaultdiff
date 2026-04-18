package vault

import (
	"context"
	"os"
	"testing"
)

func TestRollback_Integration(t *testing.T) {
	if os.Getenv("VAULT_INTEGRATION") == "" {
		t.Skip("skipping integration test; set VAULT_INTEGRATION to run")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	mount := "secret"
	path := "test/rollback"
	ctx := context.Background()

	// Write version 1
	_, err = client.Logical().WriteWithContext(ctx,
		mount+"/data/"+path,
		map[string]interface{}{"data": map[string]interface{}{"key": "v1"}},
	)
	if err != nil {
		t.Fatalf("write v1: %v", err)
	}

	// Write version 2
	_, err = client.Logical().WriteWithContext(ctx,
		mount+"/data/"+path,
		map[string]interface{}{"data": map[string]interface{}{"key": "v2"}},
	)
	if err != nil {
		t.Fatalf("write v2: %v", err)
	}

	rollbacker := NewRollbacker(client, mount)
	result, err := rollbacker.Rollback(ctx, path, 1)
	if err != nil {
		t.Fatalf("Rollback: %v", err)
	}

	if result.ToVersion != 1 {
		t.Errorf("expected ToVersion 1, got %d", result.ToVersion)
	}
	if result.FromVersion <= 2 {
		t.Errorf("expected new version > 2, got %d", result.FromVersion)
	}
}
