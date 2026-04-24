package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestAnnotator() *Annotator {
	client := &api.Client{}
	return NewAnnotator(client, "secret", "data")
}

func TestNewAnnotator_NotNil(t *testing.T) {
	a := newTestAnnotator()
	if a == nil {
		t.Fatal("expected non-nil Annotator")
	}
}

func TestAnnotator_Fields(t *testing.T) {
	a := newTestAnnotator()
	if a.client == nil {
		t.Error("expected client to be set")
	}
	if a.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", a.mount)
	}
	if a.path != "data" {
		t.Errorf("expected path 'data', got %q", a.path)
	}
}

func TestAnnotation_Fields(t *testing.T) {
	ann := Annotation{
		Version: 3,
		Author:  "alice",
		Note:    "bumped db password",
	}
	if ann.Version != 3 {
		t.Errorf("expected version 3, got %d", ann.Version)
	}
	if ann.Author != "alice" {
		t.Errorf("expected author 'alice', got %q", ann.Author)
	}
	if ann.Note != "bumped db password" {
		t.Errorf("unexpected note: %q", ann.Note)
	}
}

func TestAnnotate_NilClient(t *testing.T) {
	a := &Annotator{client: nil, mount: "secret", path: "data"}
	err := a.Annotate(1, "alice", "note")
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestGetAnnotations_NilClient(t *testing.T) {
	a := &Annotator{client: nil, mount: "secret", path: "data"}
	_, err := a.GetAnnotations()
	if err == nil {
		t.Error("expected error for nil client")
	}
}
