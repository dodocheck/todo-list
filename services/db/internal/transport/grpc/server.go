package grpc

import (
	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/dodocheck/go-pet-project-1/services/db/pb"
)

type Server struct {
	pb.UnimplementedTasksServiceServer
	dbController app.DBController
}

func NewServer(dbController app.DBController) *Server {
	return &Server{dbController: dbController}
}
