package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestTemplater() *Templater {
	client, _ := vaultapi.NewClient(vaultapi.DefaultConfig())
	return NewTemplater(client, "secret")
}

func TestNewTemplater_NotNil(t *testing.T) {
	tr := newTestTemplater()
	if tr == nil {
		t.Fatal("expected non-nil Templater")
	}
}

func TestTemplater_Fields(t *testing.T) {
	tr := newTestTemplater()
	if tr.client == nil {
		t.Error("expected client to be set")
	}
	if tr.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", tr.mount)
	}
}

func TestTemplateResult_Fields(t *testing.T) {
	res := &TemplateResult{
		Rendered:    "hello world",
		SecretsUsed: []string{"db/creds:password"},
	}
	if res.Rendered != "hello world" {
		t.Errorf("unexpected Rendered: %q", res.Rendered)
	}
	if len(res.SecretsUsed) != 1 {
		t.Errorf("expected 1 secret used, got %d", len(res.SecretsUsed))
	}
}

func TestRender_NilClient(t *testing.T) {
	tr := &Templater{client: nil, mount: "secret"}
	_, err := tr.Render(context.Background(), "{{ secret \"db/creds\" \"password\" }}")
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestParseDirective_Valid(t *testing.T) {
	path, key, err := parseDirective(`{{ secret "myapp/config" "api_key" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "myapp/config" {
		t.Errorf("expected path 'myapp/config', got %q", path)
	}
	if key != "api_key" {
		t.Errorf("expected key 'api_key', got %q", key)
	}
}

func TestParseDirective_Invalid(t *testing.T) {
	_, _, err := parseDirective(`{{ secret "badformat" }}`)
	if err == nil {
		t.Error("expected error for invalid directive")
	}
}

func TestRender_UnclosedDirective(t *testing.T) {
	tr := newTestTemplater()
	_, err := tr.Render(context.Background(), `host: {{ secret "db/host" "value"`)
	if err == nil {
		t.Error("expected error for unclosed directive")
	}
}
