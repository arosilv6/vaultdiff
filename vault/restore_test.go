package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestRestorer(t *testing.T, handler http.Handler) (*Restorer, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	cfg := api.DefaultConfig()
	cfg.Address = srv.URL
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return NewRestorer(client, "secret"), srv
}

func TestReadVersion_ReturnsData(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	})
	r, srv := newTestRestorer(t, handler)
	defer srv.Close()

	data, err := r.ReadVersion("myapp/config", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("expected value, got %v", data["key"])
	}
}

func TestReadVersion_EmptyResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r, srv := newTestRestorer(t, handler)
	defer srv.Close()

	_, err := r.ReadVersion("myapp/config", 99)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
