package vault

import (
	"context"
	"fmt"
	"strconv"

	vaultapi "github.com/hashicorp/vault/api"
)

// VersionMetadata holds metadata for a single secret version.
type VersionMetadata struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// SecretVersion represents a KV v2 secret at a specific version.
type SecretVersion struct {
	Path     string
	Version  int
	Data     map[string]string
	Metadata VersionMetadata
}

// GetSecretVersion retrieves a KV v2 secret at the given version.
// If version is 0, the latest version is returned.
func GetSecretVersion(client *vaultapi.Client, mount, secretPath string, version int) (*SecretVersion, error) {
	ctx := context.Background()

	kvClient := client.KVv2(mount)

	var secret *vaultapi.KVSecret
	var err error

	if version == 0 {
		secret, err = kvClient.Get(ctx, secretPath)
	} else {
		secret, err = kvClient.GetVersion(ctx, secretPath, version)
	}
	if err != nil {
		return nil, fmt.Errorf("reading secret %s: %w", secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %s not found", secretPath)
	}

	data := flattenCustomMetadata(secret.Data)

	meta := VersionMetadata{}
	if secret.VersionMetadata != nil {
		meta.Version = secret.VersionMetadata.Version
		meta.CreatedTime = secret.VersionMetadata.CreatedTime.String()
		if !secret.VersionMetadata.DeletionTime.IsZero() {
			meta.DeletionTime = secret.VersionMetadata.DeletionTime.String()
		}
		meta.Destroyed = secret.VersionMetadata.Destroyed
	}

	if meta.Version == 0 && version != 0 {
		meta.Version = version
	}

	return &SecretVersion{
		Path:     secretPath,
		Version:  meta.Version,
		Data:     data,
		Metadata: meta,
	}, nil
}

// VersionLabel returns a human-readable label for a version number.
func VersionLabel(version int) string {
	if version == 0 {
		return "latest"
	}
	return strconv.Itoa(version)
}
