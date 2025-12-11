package web

import (
	"errors"
	"net/http"
	"pet1/services/api"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	httpHandlers *HttpHandlers
}

func NewHttpServer(toDoList *api.ToDoList) *HttpServer {
	return &HttpServer{httpHandlers: NewHttpHandlers(toDoList)}
}

func (s *HttpServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.handleAddTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.handleListAllTasks)
	router.Path("/tasks").Methods("DELETE").HandlerFunc(s.httpHandlers.handleDeleteTask)
	router.Path("/tasks").Methods("PATCH").HandlerFunc(s.httpHandlers.handleFinishTask)

	server := http.Server{Addr: "localhost:9090", Handler: router}

	s.httpHandlers.SetCloseServerFunc(server.Close)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
