package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestExporter() *Exporter {
	return NewExporter(nil, "secret")
}

func TestNewExporter_NotNil(t *testing.T) {
	e := newTestExporter()
	if e == nil {
		t.Fatal("expected non-nil Exporter")
	}
}

func TestExporter_Fields(t *testing.T) {
	client := &api.Client{}
	e := NewExporter(client, "kv")
	if e.mount != "kv" {
		t.Errorf("expected mount kv, got %s", e.mount)
	}
	if e.client != client {
		t.Error("expected client to be set")
	}
}

func TestExportVersion_NilClient(t *testing.T) {
	e := newTestExporter()
	_, err := e.ExportVersion("myapp/config", 1)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestExportEntry_Fields(t *testing.T) {
	entry := &ExportEntry{
		Path:    "myapp/config",
		Version: 3,
		Data:    map[string]string{"key": "value"},
	}
	if entry.Path != "myapp/config" {
		t.Errorf("unexpected path: %s", entry.Path)
	}
	if entry.Version != 3 {
		t.Errorf("unexpected version: %d", entry.Version)
	}
	if entry.Data["key"] != "value" {
		t.Errorf("unexpected data value")
	}
}

func TestWriteToFile_CreatesValidJSON(t *testing.T) {
	e := newTestExporter()
	entries := []*ExportEntry{
		{Path: "a/b", Version: 1, Data: map[string]string{"foo": "bar"}},
		{Path: "c/d", Version: 2, Data: map[string]string{"baz": "qux"}},
	}

	tmp := filepath.Join(t.TempDir(), "export.json")
	if err := e.WriteToFile(entries, tmp); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	f, err := os.Open(tmp)
	if err != nil {
		t.Fatalf("opening file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var got []*ExportEntry
	for dec.More() {
		var entry ExportEntry
		if err := dec.Decode(&entry); err != nil {
			t.Fatalf("decoding entry: %v", err)
		}
		got = append(got, &entry)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Path != "a/b" || got[1].Path != "c/d" {
		t.Errorf("unexpected paths in output")
	}
}
