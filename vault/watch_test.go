package vault

import (
	"testing"
	"time"
)

func newTestWatcher() *Watcher {
	c := &Client{}
	return NewWatcher(c, "secret/data/test", time.Second)
}

func TestNewWatcher_NotNil(t *testing.T) {
	w := newTestWatcher()
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestWatcher_Fields(t *testing.T) {
	w := newTestWatcher()
	if w.path != "secret/data/test" {
		t.Errorf("unexpected path: %s", w.path)
	}
	if w.interval != time.Second {
		t.Errorf("unexpected interval: %v", w.interval)
	}
}

func TestWatchEvent_Fields(t *testing.T) {
	e := WatchEvent{
		Path:       "secret/data/foo",
		OldVersion: 1,
		NewVersion: 2,
	}
	if e.Path != "secret/data/foo" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if e.OldVersion != 1 || e.NewVersion != 2 {
		t.Errorf("unexpected versions: %d -> %d", e.OldVersion, e.NewVersion)
	}
}

func TestWatcher_NilClient(t *testing.T) {
	w := NewWatcher(nil, "secret/data/test", time.Millisecond*10)
	if w.client != nil {
		t.Error("expected nil client to be stored as-is")
	}
}
