package main

import (
	"log"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/storage/postgres"
	"github.com/dodocheck/go-pet-project-1/services/db/internal/transport/http"
	"github.com/k0kubun/pp/v3"
)

func main() {
	dbController := postgres.NewPostgresController()
	httpServer := http.NewHttpServer(dbController)

	pp.Println(dbController.ListAllTasks())

	if err := httpServer.StartServer(); err != nil {
		log.Fatal("Failed to start http web server:", err)
	}

}
