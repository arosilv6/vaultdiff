package vault

import (
	"context"
	"testing"

	hashivault "github.com/hashicorp/vault/api"
)

func newTestGarbageCollector() *GarbageCollector {
	cfg := hashivault.DefaultConfig()
	client, _ := hashivault.NewClient(cfg)
	return NewGarbageCollector(client, "secret")
}

func TestNewGarbageCollector_NotNil(t *testing.T) {
	gc := newTestGarbageCollector()
	if gc == nil {
		t.Fatal("expected non-nil GarbageCollector")
	}
}

func TestGarbageCollector_Fields(t *testing.T) {
	gc := newTestGarbageCollector()
	if gc.client == nil {
		t.Error("expected client to be set")
	}
	if gc.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", gc.mount)
	}
}

func TestGCResult_Fields(t *testing.T) {
	r := &GCResult{
		Path:           "myapp/config",
		VersionsPurged: []int{1, 2, 3},
		Count:          3,
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Count != 3 {
		t.Errorf("unexpected count: %d", r.Count)
	}
	if len(r.VersionsPurged) != 3 {
		t.Errorf("expected 3 purged versions, got %d", len(r.VersionsPurged))
	}
}

func TestCollect_NilClient(t *testing.T) {
	gc := &GarbageCollector{client: nil, mount: "secret"}
	_, err := gc.Collect(context.Background(), "myapp/config", 5)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCollect_InvalidMaxVersions(t *testing.T) {
	gc := newTestGarbageCollector()
	_, err := gc.Collect(context.Background(), "myapp/config", 0)
	if err == nil {
		t.Fatal("expected error for maxVersions < 1")
	}
}
