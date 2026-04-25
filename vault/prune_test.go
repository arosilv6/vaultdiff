package vault

import (
	"testing"

	hashivault "github.com/hashicorp/vault/api"
)

func newTestPruner() *Pruner {
	client, _ := hashivault.NewClient(hashivault.DefaultConfig())
	return NewPruner(client, "secret")
}

func TestNewPruner_NotNil(t *testing.T) {
	p := newTestPruner()
	if p == nil {
		t.Fatal("expected non-nil Pruner")
	}
}

func TestPruner_Fields(t *testing.T) {
	p := newTestPruner()
	if p.mountPath != "secret" {
		t.Errorf("expected mountPath 'secret', got %q", p.mountPath)
	}
	if p.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestPrune_NilClient(t *testing.T) {
	p := &Pruner{client: nil, mountPath: "secret"}
	_, err := p.Prune(nil, "myapp/config", 3, false) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPrune_InvalidKeep(t *testing.T) {
	p := newTestPruner()
	_, err := p.Prune(nil, "myapp/config", 0, false) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for keep < 1")
	}
}

func TestPruneResult_Fields(t *testing.T) {
	r := &PruneResult{
		Path:           "myapp/config",
		VersionsPruned: []int{1, 2, 3},
		VersionsKept:   2,
		DryRun:         true,
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if len(r.VersionsPruned) != 3 {
		t.Errorf("expected 3 pruned versions, got %d", len(r.VersionsPruned))
	}
	if r.VersionsKept != 2 {
		t.Errorf("expected 2 kept, got %d", r.VersionsKept)
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}
