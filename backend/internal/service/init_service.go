package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vinneyto/splatmaker/backend/internal/model"
	"github.com/vinneyto/splatmaker/backend/internal/octree"
	"github.com/vinneyto/splatmaker/backend/internal/ply"
)

type InitService struct {
	Client *http.Client
}

func NewInitService() *InitService {
	return &InitService{Client: &http.Client{Timeout: 30 * time.Second}}
}

func (s *InitService) BuildFromURL(ctx context.Context, plyURL string, cellSize float64, maxDepth uint32) (model.Octree, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, plyURL, nil)
	if err != nil {
		return model.Octree{}, err
	}
	resp, err := s.Client.Do(req)
	if err != nil {
		return model.Octree{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model.Octree{}, fmt.Errorf("failed to download ply: status %d", resp.StatusCode)
	}
	points, err := ply.ParseASCII(resp.Body)
	if err != nil {
		return model.Octree{}, err
	}
	return octree.BuildUniform(points, cellSize, maxDepth), nil
}
