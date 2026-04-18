package diff

import (
	"fmt"
	"sort"
)

// ChangeType describes how a key changed between versions.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string     `json:"key"`
	Type     ChangeType `json:"type"`
	OldValue string     `json:"old_value,omitempty"`
	NewValue string     `json:"new_value,omitempty"`
}

// Compare returns the list of changes between two secret maps.
func Compare(a, b map[string]string) []Change {
	var changes []Change
	for _, key := range unionKeys(a, b) {
		oldVal, inA := a[key]
		newVal, inB := b[key]
		switch {
		case inA && !inB:
			changes = append(changes, Change{Key: key, Type: Removed, OldValue: oldVal})
		case !inA && inB:
			changes = append(changes, Change{Key: key, Type: Added, NewValue: newVal})
		case oldVal != newVal:
			changes = append(changes, Change{Key: key, Type: Modified, OldValue: oldVal, NewValue: newVal})
		}
	}
	return changes
}

// Print renders changes to stdout in a human-readable format.
func Print(changes []Change) {
	if len(changes) == 0 {
		fmt.Println("No changes.")
		return
	}
	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Printf("+ %s = %s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Printf("- %s = %s\n", c.Key, c.OldValue)
		case Modified:
			fmt.Printf("~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
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
	sort.Strings(keys)
	return keys
}
