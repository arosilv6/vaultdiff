package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestCopier(t *testing.T, handler http.Handler) *Copier {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	return NewCopier(client, "secret")
}

func TestCopyVersion_Success(t *testing.T) {
	calls := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			w.Write([]byte(`{"data":{"data":{"key":"value"},"metadata":{"version":1}}}`))
			return
		}
		w.Write([]byte(`{"data":{"version":2}}`))
	})
	c := newTestCopier(t, h)
	err := c.CopyVersion(context.Background(), "src/secret", 1, "dst/secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 HTTP calls, got %d", calls)
	}
}

func TestCopyVersion_NilResponse(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	})
	c := newTestCopier(t, h)
	err := c.CopyVersion(context.Background(), "src/secret", 1, "dst/secret")
	if err == nil {
		t.Fatal("expected error for nil data")
	}
}
