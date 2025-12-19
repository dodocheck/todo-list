package http

import (
	"errors"
	"net/http"
	"os"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	httpHandlers *HttpHandlers
}

func NewHttpServer(service *app.Service) *HttpServer {
	return &HttpServer{httpHandlers: NewHttpHandlers(service)}
}

func (s *HttpServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/create").Methods("POST").HandlerFunc(s.httpHandlers.handleAddTask)
	router.Path("/list").Methods("GET").HandlerFunc(s.httpHandlers.handleListAllTasks)
	router.Path("/delete").Methods("DELETE").HandlerFunc(s.httpHandlers.handleDeleteTask)
	router.Path("/done").Methods("PUT").HandlerFunc(s.httpHandlers.handleFinishTask)

	server := http.Server{Addr: ":" + os.Getenv("API_SERVICE_INTERNAL_PORT"), Handler: router}

	s.httpHandlers.SetCloseServerFunc(server.Close)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
