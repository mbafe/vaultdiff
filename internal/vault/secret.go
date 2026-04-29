package vault

import (
	"context"
	"fmt"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretVersion represents a single version of a KV v2 secret.
type SecretVersion struct {
	Version  int
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

// SecretReader defines the interface for reading secret versions.
type SecretReader interface {
	ReadSecretVersion(ctx context.Context, mount, path string, version int) (*SecretVersion, error)
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// If version is 0, the latest version is returned.
func (c *Client) ReadSecretVersion(ctx context.Context, mount, path string, version int) (*SecretVersion, error) {
	kvClient := c.vault.KVv2(mount)

	var secret *vaultapi.KVSecret
	var err error

	if version == 0 {
		secret, err = kvClient.Get(ctx, path)
	} else {
		secret, err = kvClient.GetVersion(ctx, path, version)
	}
	if err != nil {
		return nil, fmt.Errorf("reading secret %s/%s version %d: %w", mount, path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %s/%s not found", mount, path)
	}

	v := 0
	if secret.VersionMetadata != nil {
		v = secret.VersionMetadata.Version
	}

	return &SecretVersion{
		Version:  v,
		Data:     secret.Data,
		Metadata: flattenCustomMetadata(secret.CustomMetadata),
	}, nil
}

// flattenCustomMetadata converts custom metadata values to strings.
func flattenCustomMetadata(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		out[k] = strconv.Quote(fmt.Sprintf("%v", v))
	}
	return out
}
