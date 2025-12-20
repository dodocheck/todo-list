package grpc

import (
	"log"
	"net"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedTasksServiceServer
	service *app.Service
}

func NewServer(service *app.Service) *Server {
	return &Server{service: service}
}

func (s *Server) StartServer(serverAddress string) error {
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("listen %s: %v\n", serverAddress, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTasksServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
