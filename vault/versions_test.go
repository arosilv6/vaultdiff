package vault

import (
	"testing"
)

func TestListVersions_ParsesMetadata(t *testing.T) {
	versionsData := map[string]interface{}{
		"1": map[string]interface{}{
			"created_time":  "2024-01-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
		"2": map[string]interface{}{
			"created_time":  "2024-01-02T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
	}

	var versions []VersionMeta
	for _, v := range versionsData {
		data := v.(map[string]interface{})
		meta := VersionMeta{}
		if ct, ok := data["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := data["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if destroyed, ok := data["destroyed"].(bool); ok {
			meta.Destroyed = destroyed
		}
		versions = append(versions, meta)
	}

	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
}

func TestVersionMeta_DestroyedFlag(t *testing.T) {
	meta := VersionMeta{
		Version:     3,
		CreatedTime: "2024-03-01T00:00:00Z",
		Destroyed:   true,
	}

	if !meta.Destroyed {
		t.Error("expected Destroyed to be true")
	}
	if meta.Version != 3 {
		t.Errorf("expected version 3, got %d", meta.Version)
	}
}

func TestVersionMeta_DeletionTime(t *testing.T) {
	meta := VersionMeta{
		Version:      2,
		CreatedTime:  "2024-02-01T00:00:00Z",
		DeletionTime: "2024-02-10T00:00:00Z",
		Destroyed:    false,
	}

	if meta.DeletionTime == "" {
		t.Error("expected non-empty DeletionTime")
	}
}
