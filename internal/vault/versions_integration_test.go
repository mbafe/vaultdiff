//go:build integration
// +build integration

package vault_test

import (
	"os"
	"testing"

	"github.com/user/vaultdiff/internal/vault"
)

// TestGetSecretVersion_Integration requires a running Vault instance.
// Run with: go test -tags=integration ./internal/vault/...
//
// Required environment variables:
//   VAULT_ADDR  - Vault server address (default: http://127.0.0.1:8200)
//   VAULT_TOKEN - Vault token
//   VAULT_TEST_PATH - KV v2 path to read (default: secret/test)
func TestGetSecretVersion_Integration(t *testing.T) {
	client, err := vault.NewClient("", "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	mount := envOrDefault("VAULT_TEST_MOUNT", "secret")
	path := envOrDefault("VAULT_TEST_PATH", "test")

	sv, err := vault.GetSecretVersion(client, mount, path, 0)
	if err != nil {
		t.Fatalf("GetSecretVersion: %v", err)
	}

	if sv.Path != path {
		t.Errorf("expected path %q, got %q", path, sv.Path)
	}
	if sv.Version < 1 {
		t.Errorf("expected version >= 1, got %d", sv.Version)
	}
	t.Logf("secret %s version %s has %d keys",
		sv.Path, vault.VersionLabel(sv.Version), len(sv.Data))
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
