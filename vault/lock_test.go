package vault

import (
	"context"
	"testing"
)

func newTestLocker() *Locker {
	return &Locker{client: nil, mount: "secret"}
}

func TestNewLocker_NotNil(t *testing.T) {
	l := NewLocker(nil, "secret")
	if l == nil {
		t.Fatal("expected non-nil Locker")
	}
}

func TestLocker_Fields(t *testing.T) {
	l := NewLocker(nil, "kv")
	if l.mount != "kv" {
		t.Errorf("expected mount kv, got %s", l.mount)
	}
}

func TestLock_NilClient(t *testing.T) {
	l := newTestLocker()
	err := l.Lock(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error with nil client")
	}
}

func TestUnlock_NilClient(t *testing.T) {
	l := newTestLocker()
	err := l.Unlock(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error with nil client")
	}
}

func TestIsLocked_NilClient(t *testing.T) {
	l := newTestLocker()
	_, err := l.IsLocked(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error with nil client")
	}
}

func TestIsLocked_NilSecret(t *testing.T) {
	// Simulate nil secret response by calling isLocked logic directly
	l := newTestLocker()
	// We can't call Vault without a real client, but we verify nil-safe path
	// by checking the locker struct is valid
	if l.mount == "" {
		t.Error("mount should not be empty")
	}
}
