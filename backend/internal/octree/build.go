package octree

import (
	"math"

	"github.com/vinneyto/splatmaker/backend/internal/model"
)

func BuildUniform(points []model.Vec3, cellSize float64, maxDepth uint32) model.Octree {
	if len(points) == 0 {
		return model.Octree{}
	}
	if cellSize <= 0 {
		cellSize = 1.0
	}
	min, max := bounds(points)
	root := model.OctreeNode{ID: 0, Depth: 0, Min: min, Max: max}
	nodes := []model.OctreeNode{root}
	indices := make([]uint32, len(points))
	for i := range points {
		indices[i] = uint32(i)
	}
	split(&nodes, 0, points, indices, cellSize, maxDepth)
	return model.Octree{Min: min, Max: max, Nodes: nodes}
}

func split(nodes *[]model.OctreeNode, nodeID uint32, points []model.Vec3, splatIndices []uint32, cellSize float64, maxDepth uint32) {
	n := &(*nodes)[nodeID]
	sx := n.Max.X - n.Min.X
	sy := n.Max.Y - n.Min.Y
	sz := n.Max.Z - n.Min.Z
	if n.Depth >= maxDepth || (sx <= cellSize && sy <= cellSize && sz <= cellSize) || len(splatIndices) <= 1 {
		n.SplatIndices = append(n.SplatIndices, splatIndices...)
		return
	}

	mx := (n.Min.X + n.Max.X) * 0.5
	my := (n.Min.Y + n.Max.Y) * 0.5
	mz := (n.Min.Z + n.Max.Z) * 0.5
	buckets := make([][]uint32, 8)
	for _, si := range splatIndices {
		p := points[int(si)]
		i := 0
		if p.X >= mx { i |= 1 }
		if p.Y >= my { i |= 2 }
		if p.Z >= mz { i |= 4 }
		buckets[i] = append(buckets[i], si)
	}

	for i, bucket := range buckets {
		if len(bucket) == 0 {
			continue
		}
		child := model.OctreeNode{ID: uint32(len(*nodes)), Depth: n.Depth + 1}
		child.Min = model.Vec3{X: pick(n.Min.X, mx, i&1 != 0), Y: pick(n.Min.Y, my, i&2 != 0), Z: pick(n.Min.Z, mz, i&4 != 0)}
		child.Max = model.Vec3{X: pick(mx, n.Max.X, i&1 != 0), Y: pick(my, n.Max.Y, i&2 != 0), Z: pick(mz, n.Max.Z, i&4 != 0)}
		*nodes = append(*nodes, child)
		n.ChildrenIDs = append(n.ChildrenIDs, child.ID)
	}

	for i, bucket := range buckets {
		if len(bucket) == 0 { continue }
		childID := n.ChildrenIDs[indexInNonEmpty(i, buckets)]
		split(nodes, childID, points, bucket, cellSize, maxDepth)
	}
}

func indexInNonEmpty(i int, buckets [][]uint32) int {
	j := 0
	for k:=0; k<i; k++ { if len(buckets[k])>0 { j++ } }
	return j
}

func pick(a, b float64, upper bool) float64 {
	if upper { return b }
	return a
}

func bounds(points []model.Vec3) (model.Vec3, model.Vec3) {
	min := model.Vec3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	max := model.Vec3{X: -math.MaxFloat64, Y: -math.MaxFloat64, Z: -math.MaxFloat64}
	for _, p := range points {
		if p.X < min.X { min.X = p.X }
		if p.Y < min.Y { min.Y = p.Y }
		if p.Z < min.Z { min.Z = p.Z }
		if p.X > max.X { max.X = p.X }
		if p.Y > max.Y { max.Y = p.Y }
		if p.Z > max.Z { max.Z = p.Z }
	}
	return min, max
}
