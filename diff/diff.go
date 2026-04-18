package diff

import (
	"fmt"
	"io"
	"sort"
)

// ChangeType describes the kind of change for a key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level diff entry.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare returns the diff between two secret maps.
func Compare(a, b map[string]string) []Change {
	keys := unionKeys(a, b)
	sort.Strings(keys)

	var changes []Change
	for _, k := range keys {
		oldVal, inA := a[k]
		newVal, inB := b[k]

		switch {
		case inA && !inB:
			changes = append(changes, Change{Key: k, Type: Removed, OldValue: oldVal})
		case !inA && inB:
			changes = append(changes, Change{Key: k, Type: Added, NewValue: newVal})
		case oldVal != newVal:
			changes = append(changes, Change{Key: k, Type: Modified, OldValue: oldVal, NewValue: newVal})
		default:
			changes = append(changes, Change{Key: k, Type: Unchanged})
		}
	}
	return changes
}

// Print writes a human-readable diff to w.
func Print(w io.Writer, changes []Change) {
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "+ %s = %s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(w, "- %s = %s\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(w, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
