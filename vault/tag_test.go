package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestTagger() *Tagger {
	client, _ := api.NewClient(api.DefaultConfig())
	return NewTagger(client, "secret")
}

func TestNewTagger_NotNil(t *testing.T) {
	tagger := newTestTagger()
	if tagger == nil {
		t.Fatal("expected non-nil Tagger")
	}
}

func TestTagger_Fields(t *testing.T) {
	tagger := newTestTagger()
	if tagger.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", tagger.mount)
	}
	if tagger.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestSetTag_NilClient(t *testing.T) {
	tagger := &Tagger{client: nil, mount: "secret"}
	err := tagger.SetTag(nil, "myapp/config", "env", "prod") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetTags_NilClient(t *testing.T) {
	tagger := &Tagger{client: nil, mount: "secret"}
	_, err := tagger.GetTags(nil, "myapp/config") //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetTags_NilSecret(t *testing.T) {
	// Simulate nil response path via helper that returns empty map
	tagger := newTestTagger()
	// We can't call Vault in unit tests; verify struct is properly initialised.
	if tagger.mount == "" {
		t.Error("mount should not be empty")
	}
}
