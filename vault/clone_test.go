package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestCloner() *Cloner {
	cfg := api.DefaultConfig()
	client, _ := api.NewClient(cfg)
	return NewCloner(client, "secret")
}

func TestNewCloner_NotNil(t *testing.T) {
	c := newTestCloner()
	if c == nil {
		t.Fatal("expected non-nil Cloner")
	}
}

func TestCloner_Fields(t *testing.T) {
	c := newTestCloner()
	if c.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", c.mount)
	}
	if c.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestClone_NilClient(t *testing.T) {
	c := &Cloner{client: nil, mount: "secret"}
	_, err := c.Clone("src/key", "dst/key")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCloneResult_Fields(t *testing.T) {
	r := &CloneResult{
		SourcePath: "src/key",
		DestPath:   "dst/key",
		Versions:   []int{1, 2, 3},
	}
	if r.SourcePath != "src/key" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "dst/key" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if len(r.Versions) != 3 {
		t.Errorf("expected 3 versions, got %d", len(r.Versions))
	}
}

func TestCloneResult_EmptyVersions(t *testing.T) {
	r := &CloneResult{
		SourcePath: "src/empty",
		DestPath:   "dst/empty",
		Versions:   []int{},
	}
	if len(r.Versions) != 0 {
		t.Errorf("expected 0 versions, got %d", len(r.Versions))
	}
}
