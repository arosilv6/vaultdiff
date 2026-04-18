package vault

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestWatch_Integration(t *testing.T) {
	if os.Getenv("VAULT_ADDR") == "" {
		t.Skip("skipping integration test: VAULT_ADDR not set")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	path := os.Getenv("VAULT_TEST_PATH")
	if path == "" {
		path = "secret/data/vaultdiff-watch-test"
	}

	watcher := NewWatcher(client, path, 500*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := watcher.Watch(ctx)

	// Drain events until context expires.
	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return
			}
			if ev.Error != nil {
				// Errors are acceptable in integration env without a real secret.
				t.Logf("watch error (expected in empty env): %v", ev.Error)
			}
		case <-ctx.Done():
			return
		}
	}
}
