package vault

import (
	"testing"

	hashivault "github.com/hashicorp/vault/api"
)

func newTestSummarizer() *Summarizer {
	cfg := hashivault.DefaultConfig()
	cfg.Address = "http://127.0.0.1:8200"
	client, _ := hashivault.NewClient(cfg)
	return NewSummarizer(client, "secret")
}

func TestNewSummarizer_NotNil(t *testing.T) {
	s := newTestSummarizer()
	if s == nil {
		t.Fatal("expected non-nil Summarizer")
	}
}

func TestSummarizer_Fields(t *testing.T) {
	s := newTestSummarizer()
	if s.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", s.mount)
	}
	if s.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestSummarize_NilClient(t *testing.T) {
	s := &Summarizer{client: nil, mount: "secret"}
	_, err := s.Summarize("myapp/config")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestSummaryEntry_Fields(t *testing.T) {
	entry := &SummaryEntry{
		Path:           "myapp/config",
		VersionCount:   5,
		LatestVersion:  5,
		DeletedCount:   1,
		DestroyedCount: 0,
	}
	if entry.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", entry.Path)
	}
	if entry.VersionCount != 5 {
		t.Errorf("unexpected VersionCount: %d", entry.VersionCount)
	}
	if entry.LatestVersion != 5 {
		t.Errorf("unexpected LatestVersion: %d", entry.LatestVersion)
	}
	if entry.DeletedCount != 1 {
		t.Errorf("unexpected DeletedCount: %d", entry.DeletedCount)
	}
	if entry.DestroyedCount != 0 {
		t.Errorf("unexpected DestroyedCount: %d", entry.DestroyedCount)
	}
}

func TestNewSummarizer_CustomMount(t *testing.T) {
	cfg := hashivault.DefaultConfig()
	client, _ := hashivault.NewClient(cfg)
	s := NewSummarizer(client, "kv")
	if s.mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", s.mount)
	}
}
