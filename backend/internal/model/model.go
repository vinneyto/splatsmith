package model

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

type Octree struct {
	Min   Vec3
	Max   Vec3
	Nodes []OctreeNode
}

type OctreeNode struct {
	ID           uint32
	Depth        uint32
	Min          Vec3
	Max          Vec3
	ChildrenIDs  []uint32
	SplatIndices []uint32
}
