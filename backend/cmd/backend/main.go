package main

import (
	"log"
	"net"

	backendv1 "github.com/vinneyto/splatmaker/backend/gen/backend/v1"
	"github.com/vinneyto/splatmaker/backend/internal/grpcserver"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	backendv1.RegisterSplatIndexServiceServer(srv, grpcserver.NewSplatIndexServer())
	log.Printf("backend grpc listening on %s", lis.Addr().String())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
