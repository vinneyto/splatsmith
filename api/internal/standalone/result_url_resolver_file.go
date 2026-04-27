package standalone

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vinneyto/splatra/api/internal/core"
)

type FileResultURLResolver struct {
	root string
}

func NewFileResultURLResolver(root string) (*FileResultURLResolver, error) {
	if root == "" {
		return nil, fmt.Errorf("results root is empty")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absRoot, 0o755); err != nil {
		return nil, err
	}
	return &FileResultURLResolver{root: absRoot}, nil
}

func (r *FileResultURLResolver) ResolveResultURL(_ context.Context, key string, ttl time.Duration) (core.ResultFileURL, error) {
	if key == "" {
		return core.ResultFileURL{}, fmt.Errorf("result key is empty")
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	cleanKey := filepath.Clean(strings.TrimLeft(key, "/"))
	absPath := filepath.Join(r.root, cleanKey)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return core.ResultFileURL{}, err
	}

	return core.ResultFileURL{
		Key:       key,
		URL:       "file://" + absPath,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}, nil
}
