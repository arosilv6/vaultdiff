package vault

import (
	"testing"
	"time"
)

func newTestHistorian() *Historian {
	return NewHistorian(nil, "secret")
}

func TestNewHistorian_NotNil(t *testing.T) {
	h := newTestHistorian()
	if h == nil {
		t.Fatal("expected non-nil Historian")
	}
}

func TestHistorian_Fields(t *testing.T) {
	h := newTestHistorian()
	if h.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", h.mount)
	}
	if h.client != nil {
		t.Error("expected nil client for test historian")
	}
}

func TestHistory_NilClient(t *testing.T) {
	h := newTestHistorian()
	_, err := h.History("myapp/config")
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestHistoryEntry_Fields(t *testing.T) {
	now := time.Now()
	entry := HistoryEntry{
		Version:     3,
		CreatedTime: now,
		Deleted:     false,
		Destroyed:   false,
	}
	if entry.Version != 3 {
		t.Errorf("expected version 3, got %d", entry.Version)
	}
	if entry.CreatedTime != now {
		t.Errorf("unexpected created time")
	}
	if entry.Deleted {
		t.Error("expected Deleted to be false")
	}
	if entry.Destroyed {
		t.Error("expected Destroyed to be false")
	}
}

func TestHistoryEntry_DeletedAndDestroyed(t *testing.T) {
	entry := HistoryEntry{
		Version:   2,
		Deleted:   true,
		Destroyed: true,
	}
	if !entry.Deleted {
		t.Error("expected Deleted to be true")
	}
	if !entry.Destroyed {
		t.Error("expected Destroyed to be true")
	}
}
