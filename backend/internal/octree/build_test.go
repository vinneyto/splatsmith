package octree

import (
	"testing"

	"github.com/vinneyto/splatmaker/backend/internal/model"
)

func TestBuildUniform(t *testing.T) {
	points := []model.Vec3{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}, {3, 3, 3}}
	tree := BuildUniform(points, 1.0, 4)
	if len(tree.Nodes) == 0 {
		t.Fatalf("expected non-empty tree")
	}
	if tree.Nodes[0].ID != 0 {
		t.Fatalf("root id must be 0")
	}
}
