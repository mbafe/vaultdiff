package vault

import (
	"testing"
)

func TestVersionLabel_Zero(t *testing.T) {
	label := VersionLabel(0)
	if label != "latest" {
		t.Errorf("expected \"latest\", got %q", label)
	}
}

func TestVersionLabel_NonZero(t *testing.T) {
	cases := []struct {
		version  int
		expected string
	}{
		{1, "1"},
		{5, "5"},
		{42, "42"},
	}
	for _, tc := range cases {
		t.Run(tc.expected, func(t *testing.T) {
			got := VersionLabel(tc.version)
			if got != tc.expected {
				t.Errorf("VersionLabel(%d): expected %q, got %q", tc.version, tc.expected, got)
			}
		})
	}
}

func TestVersionMetadata_Fields(t *testing.T) {
	meta := VersionMetadata{
		Version:      3,
		CreatedTime:  "2024-01-15T10:00:00Z",
		DeletionTime: "",
		Destroyed:    false,
	}

	if meta.Version != 3 {
		t.Errorf("expected Version=3, got %d", meta.Version)
	}
	if meta.CreatedTime != "2024-01-15T10:00:00Z" {
		t.Errorf("unexpected CreatedTime: %s", meta.CreatedTime)
	}
	if meta.Destroyed {
		t.Error("expected Destroyed=false")
	}
}

func TestSecretVersionStruct(t *testing.T) {
	sv := &SecretVersion{
		Path:    "secret/myapp/config",
		Version: 2,
		Data:    map[string]string{"key": "value"},
		Metadata: VersionMetadata{
			Version:     2,
			CreatedTime: "2024-03-01T00:00:00Z",
		},
	}

	if sv.Path != "secret/myapp/config" {
		t.Errorf("unexpected Path: %s", sv.Path)
	}
	if sv.Version != 2 {
		t.Errorf("expected Version=2, got %d", sv.Version)
	}
	if sv.Data["key"] != "value" {
		t.Errorf("unexpected Data value: %s", sv.Data["key"])
	}
}
