package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func newTestComparer() *Comparer {
	client, _ := api.NewClient(api.DefaultConfig())
	return NewComparer(client, "secret")
}

func TestNewComparer_NotNil(t *testing.T) {
	c := newTestComparer()
	if c == nil {
		t.Fatal("expected non-nil Comparer")
	}
}

func TestComparer_Fields(t *testing.T) {
	c := newTestComparer()
	if c.mount != "secret" {
		t.Errorf("expected mount 'secret', got %q", c.mount)
	}
	if c.client == nil {
		t.Error("expected non-nil client")
	}
}

func TestFetchVersions_NilClient(t *testing.T) {
	c := &Comparer{client: nil, mount: "secret"}
	_, err := c.FetchVersions("myapp/config", 1, 2)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestCompareResult_Fields(t *testing.T) {
	r := &CompareResult{
		Path:     "myapp/config",
		VersionA: 1,
		VersionB: 2,
		DataA:    map[string]interface{}{"key": "old"},
		DataB:    map[string]interface{}{"key": "new"},
	}
	if r.Path != "myapp/config" {
		t.Errorf("unexpected Path: %s", r.Path)
	}
	if r.VersionA != 1 || r.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", r.VersionA, r.VersionB)
	}
	if r.DataA["key"] != "old" || r.DataB["key"] != "new" {
		t.Error("unexpected data values")
	}
}

func TestCompareResult_EmptyData(t *testing.T) {
	r := &CompareResult{
		Path:     "myapp/config",
		VersionA: 1,
		VersionB: 2,
		DataA:    map[string]interface{}{},
		DataB:    map[string]interface{}{},
	}
	if len(r.DataA) != 0 || len(r.DataB) != 0 {
		t.Error("expected empty data maps")
	}
}
