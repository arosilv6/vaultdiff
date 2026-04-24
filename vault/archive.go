package vault

import (
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// ArchiveEntry represents a single archived secret version.
type ArchiveEntry struct {
	Path      string
	Version   int
	Data      map[string]interface{}
	ArchivedAt time.Time
}

// Archiver handles archiving secret versions to a designated archive path.
type Archiver struct {
	client *vaultapi.Client
	mount  string
}

// NewArchiver creates a new Archiver instance.
func NewArchiver(client *vaultapi.Client, mount string) *Archiver {
	return &Archiver{
		client: client,
		mount:  mount,
	}
}

// Archive reads a specific version of a secret and writes it to an archive path.
func (a *Archiver) Archive(srcPath string, version int, archivePath string) (*ArchiveEntry, error) {
	if a.client == nil {
		return nil, fmt.Errorf("vault client is nil")
	}

	versionedPath := fmt.Sprintf("%s/data/%s", a.mount, srcPath)
	secret, err := a.client.Logical().ReadWithData(versionedPath, map[string][]string{
		"version": {fmt.Sprintf("%d", version)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read version %d of %s: %w", version, srcPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found for %s at version %d", srcPath, version)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for %s", srcPath)
	}

	destPath := fmt.Sprintf("%s/data/%s", a.mount, archivePath)
	_, err = a.client.Logical().Write(destPath, map[string]interface{}{"data": data})
	if err != nil {
		return nil, fmt.Errorf("failed to write archive to %s: %w", archivePath, err)
	}

	return &ArchiveEntry{
		Path:       archivePath,
		Version:    version,
		Data:       data,
		ArchivedAt: time.Now().UTC(),
	}, nil
}
