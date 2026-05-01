package grpcserver

import (
	"context"

	backendv1 "github.com/vinneyto/splatmaker/backend/gen/backend/v1"
	"github.com/vinneyto/splatmaker/backend/internal/service"
)

type SplatIndexServer struct {
	backendv1.UnimplementedSplatIndexServiceServer
	init *service.InitService
}

func NewSplatIndexServer() *SplatIndexServer {
	return &SplatIndexServer{init: service.NewInitService()}
}

func (s *SplatIndexServer) InitFromPlyUrl(ctx context.Context, req *backendv1.InitFromPlyUrlRequest) (*backendv1.InitFromPlyUrlResponse, error) {
	tree, err := s.init.BuildFromURL(ctx, req.GetPlyUrl(), req.GetCellSize(), req.GetMaxDepth())
	if err != nil {
		return nil, err
	}
	resp := &backendv1.InitFromPlyUrlResponse{Tree: &backendv1.Octree{Min: toPBVec3(tree.Min), Max: toPBVec3(tree.Max)}}
	resp.Tree.Nodes = make([]*backendv1.OctreeNode, 0, len(tree.Nodes))
	for _, n := range tree.Nodes {
		node := &backendv1.OctreeNode{
			Id:           n.ID,
			Depth:        n.Depth,
			Min:          toPBVec3(n.Min),
			Max:          toPBVec3(n.Max),
			ChildrenIds:  append([]uint32(nil), n.ChildrenIDs...),
			SplatIndices: append([]uint32(nil), n.SplatIndices...),
		}
		resp.Tree.Nodes = append(resp.Tree.Nodes, node)
	}
	return resp, nil
}
