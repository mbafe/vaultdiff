package vault

import (
	"testing"
)

func TestFlattenCustomMetadata_Nil(t *testing.T) {
	result := flattenCustomMetadata(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestFlattenCustomMetadata_Values(t *testing.T) {
	input := map[string]interface{}{
		"owner": "alice",
		"version": 3,
	}
	result := flattenCustomMetadata(input)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	for k := range input {
		if _, ok := result[k]; !ok {
			t.Errorf("expected key %q in result", k)
		}
	}
}

func TestSecretVersion_Fields(t *testing.T) {
	sv := &SecretVersion{
		Version: 2,
		Data: map[string]interface{}{
			"password": "s3cr3t",
		},
		Metadata: map[string]interface{}{
			"owner": `"alice"`,
		},
	}
	if sv.Version != 2 {
		t.Errorf("expected version 2, got %d", sv.Version)
	}
	if sv.Data["password"] != "s3cr3t" {
		t.Errorf("unexpected data value")
	}
	if sv.Metadata["owner"] != `"alice"` {
		t.Errorf("unexpected metadata value")
	}
}
