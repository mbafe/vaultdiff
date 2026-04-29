package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	vc *vaultapi.Client
}

// NewClient creates a new Vault client using environment variables or provided address/token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()

	if address != "" {
		cfg.Address = address
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	} else {
		cfg.Address = "http://127.0.0.1:8200"
	}

	if err := cfg.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("reading vault environment: %w", err)
	}

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	if token != "" {
		vc.SetToken(token)
	} else if t := os.Getenv("VAULT_TOKEN"); t != "" {
		vc.SetToken(t)
	}

	return &Client{vc: vc}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// path should be the secret path without the "secret/data/" prefix.
func (c *Client) ReadSecretVersion(mount, path string, version int) (map[string]interface{}, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.vc.KVv2(mount).GetVersion(nil, path, version)
	if err != nil {
		return nil, fmt.Errorf("reading secret %s@%d: %w", path, version, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %s version %d not found", path, version)
	}

	return secret.Data, nil
}
