package standalone

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalObjectStorage struct {
	root string
}

func NewLocalObjectStorage(root string) (*LocalObjectStorage, error) {
	if root == "" {
		return nil, fmt.Errorf("storage root is empty")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &LocalObjectStorage{root: root}, nil
}

func (s *LocalObjectStorage) SaveInputVideo(_ context.Context, userID, scanID string, data io.Reader) (string, error) {
	dir := filepath.Join(s.root, userID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, scanID+".bin")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return "", err
	}
	return path, nil
}

func (s *LocalObjectStorage) OpenResultAsset(_ context.Context, assetPath string) (io.ReadCloser, error) {
	return os.Open(assetPath)
}
