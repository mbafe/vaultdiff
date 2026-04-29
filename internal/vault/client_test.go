package vault

import (
	"os"
	"testing"
)

func TestNewClient_DefaultAddress(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	client, err := NewClient("", "test-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.vc.Token() != "test-token" {
		t.Errorf("expected token 'test-token', got '%s'", client.vc.Token())
	}
	if client.vc.Address() != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got '%s'", client.vc.Address())
	}
}

func TestNewClient_CustomAddress(t *testing.T) {
	client, err := NewClient("http://vault.example.com:8200", "my-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client.vc.Address() != "http://vault.example.com:8200" {
		t.Errorf("expected custom address, got '%s'", client.vc.Address())
	}
}

func TestNewClient_EnvVars(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://env-vault:8200")
	t.Setenv("VAULT_TOKEN", "env-token")

	client, err := NewClient("", "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client.vc.Address() != "http://env-vault:8200" {
		t.Errorf("expected env address, got '%s'", client.vc.Address())
	}
	if client.vc.Token() != "env-token" {
		t.Errorf("expected env token, got '%s'", client.vc.Token())
	}
}

func TestNewClient_ExplicitTokenOverridesEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")

	client, err := NewClient("", "explicit-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client.vc.Token() != "explicit-token" {
		t.Errorf("expected explicit token to override env, got '%s'", client.vc.Token())
	}
}
