package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompare_AddedKey(t *testing.T) {
	a := map[string]string{"foo": "bar"}
	b := map[string]string{"foo": "bar", "baz": "qux"}

	changes := Compare(a, b)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	found := findChange(changes, "baz")
	if found == nil || found.Type != Added {
		t.Errorf("expected baz to be Added")
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	a := map[string]string{"foo": "bar", "old": "val"}
	b := map[string]string{"foo": "bar"}

	changes := Compare(a, b)
	found := findChange(changes, "old")
	if found == nil || found.Type != Removed {
		t.Errorf("expected old to be Removed")
	}
}

func TestCompare_ModifiedKey(t *testing.T) {
	a := map[string]string{"key": "v1"}
	b := map[string]string{"key": "v2"}

	changes := Compare(a, b)
	found := findChange(changes, "key")
	if found == nil || found.Type != Modified {
		t.Errorf("expected key to be Modified")
	}
	if found.OldValue != "v1" || found.NewValue != "v2" {
		t.Errorf("unexpected values: %+v", found)
	}
}

func TestPrint_Output(t *testing.T) {
	changes := []Change{
		{Key: "added", Type: Added, NewValue: "new"},
		{Key: "removed", Type: Removed, OldValue: "old"},
		{Key: "modified", Type: Modified, OldValue: "a", NewValue: "b"},
		{Key: "same", Type: Unchanged},
	}

	var buf bytes.Buffer
	Print(&buf, changes)
	out := buf.String()

	if !strings.Contains(out, "+ added") {
		t.Errorf("missing added line")
	}
	if !strings.Contains(out, "- removed") {
		t.Errorf("missing removed line")
	}
	if !strings.Contains(out, "~ modified") {
		t.Errorf("missing modified line")
	}
	if strings.Contains(out, "same") {
		t.Errorf("unchanged keys should not appear in output")
	}
}

func findChange(changes []Change, key string) *Change {
	for i := range changes {
		if changes[i].Key == key {
			return &changes[i]
		}
	}
	return nil
}
