package vault

import (
	"context"
	"fmt"
	"strings"
)

// SearchResult holds a matching secret path and version metadata.
type SearchResult struct {
	Path    string
	Version int
	Meta    VersionMeta
}

// Searcher searches for secrets containing a given key or value pattern.
type Searcher struct {
	client *Client
}

// NewSearcher creates a new Searcher.
func NewSearcher(c *Client) *Searcher {
	return &Searcher{client: c}
}

// FindByKey returns all secret paths whose latest version contains the given key.
func (s *Searcher) FindByKey(ctx context.Context, mountPath, pattern string) ([]SearchResult, error) {
	paths, err := s.listRecursive(ctx, mountPath, "")
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, p := range paths {
		versions, err := ListVersions(ctx, s.client, mountPath, p)
		if err != nil || len(versions) == 0 {
			continue
		}
		latest := versions[len(versions)-1]
		data, err := s.client.ReadSecretVersion(ctx, mountPath, p, latest.Version)
		if err != nil {
			continue
		}
		for k := range data {
			if strings.Contains(k, pattern) {
				results = append(results, SearchResult{Path: p, Version: latest.Version, Meta: latest})
				break
			}
		}
	}
	return results, nil
}

// listRecursive lists all secret paths under prefix recursively.
func (s *Searcher) listRecursive(ctx context.Context, mount, prefix string) ([]string, error) {
	listPath := fmt.Sprintf("%s/metadata/%s", mount, prefix)
	secret, err := s.client.Logical().ListWithContext(ctx, listPath)
	if err != nil || secret == nil {
		return nil, err
	}
	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, nil
	}
	var paths []string
	for _, k := range keys {
		key := fmt.Sprintf("%s%s", prefix, k.(string))
		if strings.HasSuffix(key, "/") {
			sub, err := s.listRecursive(ctx, mount, key)
			if err == nil {
				paths = append(paths, sub...)
			}
		} else {
			paths = append(paths, key)
		}
	}
	return paths, nil
}
