package http

import (
	"errors"
	"net/http"
	"os"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/app"
	"github.com/gorilla/mux"
)

type Server struct {
	httpHandlers *HttpHandlers
}

func NewServer(dbController app.DBController) *Server {
	return &Server{httpHandlers: NewHttpHandlers(dbController)}
}

func (s *Server) StartServer() error {
	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.handleAddTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.handleListAllTasks)
	router.Path("/tasks").Methods("DELETE").HandlerFunc(s.httpHandlers.handleDeleteTask)
	router.Path("/tasks").Methods("PATCH").HandlerFunc(s.httpHandlers.handleFinishTask)

	server := http.Server{Addr: ":" + os.Getenv("DB_SERVICE_INTERNAL_PORT"), Handler: router}

	s.httpHandlers.SetCloseServerFunc(server.Close)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
