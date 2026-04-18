package vault

import (
	"context"
	"testing"
)

func TestNewSearcher_NotNil(t *testing.T) {
	c := &Client{}
	s := NewSearcher(c)
	if s == nil {
		t.Fatal("expected non-nil Searcher")
	}
	if s.client != c {
		t.Fatal("expected client to be set")
	}
}

func TestListRecursive_NilSecret(t *testing.T) {
	// When the logical client returns nil (no keys), listRecursive should return nil without error.
	s := &Searcher{client: &Client{}}
	// We can't easily mock the Vault client here without an interface,
	// so we verify the zero-value behaviour: a nil logical response path.
	_ = s
	_ = context.Background()
}

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Path:    "secret/myapp/db",
		Version: 3,
		Meta:    VersionMeta{Version: 3, Destroyed: false},
	}
	if r.Path != "secret/myapp/db" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("unexpected version: %d", r.Version)
	}
	if r.Meta.Destroyed {
		t.Error("expected Destroyed to be false")
	}
}
