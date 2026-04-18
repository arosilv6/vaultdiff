package vault

import (
	"testing"
	"time"
)

func newTestExpirer() *Expirer {
	return NewExpirer(nil, "secret")
}

func TestNewExpirer_NotNil(t *testing.T) {
	e := newTestExpirer()
	if e == nil {
		t.Fatal("expected non-nil Expirer")
	}
}

func TestExpirer_Fields(t *testing.T) {
	e := NewExpirer(nil, "kv")
	if e.mount != "kv" {
		t.Errorf("expected mount kv, got %s", e.mount)
	}
}

func TestCheckExpiry_NilClient(t *testing.T) {
	e := newTestExpirer()
	_, err := e.CheckExpiry(nil, "myapp/config", 1) //nolint:staticcheck
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestExpiryInfo_Fields(t *testing.T) {
	now := time.Now().Add(10 * time.Minute)
	info := &ExpiryInfo{
		Path:      "app/secret",
		Version:   3,
		ExpiresAt: now,
		TTL:       10 * time.Minute,
		Expired:   false,
	}
	if info.Path != "app/secret" {
		t.Errorf("unexpected path: %s", info.Path)
	}
	if info.Version != 3 {
		t.Errorf("unexpected version: %d", info.Version)
	}
	if info.Expired {
		t.Error("expected not expired")
	}
}

func TestExpiryInfo_ExpiredFlag(t *testing.T) {
	info := &ExpiryInfo{
		ExpiresAt: time.Now().Add(-1 * time.Minute),
		Expired:   true,
	}
	if !info.Expired {
		t.Error("expected expired to be true")
	}
}
