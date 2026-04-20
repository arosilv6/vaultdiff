package vault

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestImporter() *Importer {
	client, _ := api.NewClient(api.DefaultConfig())
	return NewImporter(client, "secret")
}

func TestNewImporter_NotNil(t *testing.T) {
	im := newTestImporter()
	if im == nil {
		t.Fatal("expected non-nil Importer")
	}
}

func TestImporter_Fields(t *testing.T) {
	im := newTestImporter()
	if im.mount != "secret" {
		t.Errorf("expected mount 'secret', got %s", im.mount)
	}
	if im.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestImportFile_NilClient(t *testing.T) {
	im := NewImporter(nil, "secret")

	entries := []ImportEntry{
		{Path: "myapp/config", Data: map[string]interface{}{"key": "value"}},
	}
	f, err := os.CreateTemp("", "import-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	json.NewEncoder(f).Encode(entries)
	f.Close()

	result, err := im.ImportFile(context.Background(), f.Name(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", result.Skipped)
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error recorded")
	}
}

func TestImportFile_DryRun(t *testing.T) {
	im := newTestImporter()

	entries := []ImportEntry{
		{Path: "app/db", Data: map[string]interface{}{"password": "s3cr3t"}},
		{Path: "app/api", Data: map[string]interface{}{"token": "abc123"}},
	}
	f, err := os.CreateTemp("", "import-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	json.NewEncoder(f).Encode(entries)
	f.Close()

	result, err := im.ImportFile(context.Background(), f.Name(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Imported != 2 {
		t.Errorf("expected 2 imported in dry-run, got %d", result.Imported)
	}
}

func TestImportResult_Fields(t *testing.T) {
	r := &ImportResult{Imported: 3, Skipped: 1, Errors: []string{"oops"}}
	if r.Imported != 3 {
		t.Errorf("expected Imported=3")
	}
	if r.Skipped != 1 {
		t.Errorf("expected Skipped=1")
	}
	if len(r.Errors) != 1 {
		t.Errorf("expected 1 error")
	}
}
