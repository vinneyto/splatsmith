package grpcserver

import (
	backendv1 "github.com/vinneyto/splatmaker/backend/gen/backend/v1"
	"github.com/vinneyto/splatmaker/backend/internal/model"
)

func toPBVec3(v model.Vec3) *backendv1.Vec3 {
	return &backendv1.Vec3{X: float32(v.X), Y: float32(v.Y), Z: float32(v.Z)}
}
