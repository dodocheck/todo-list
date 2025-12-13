package main

import (
	"log"

	"github.com/dodocheck/go-pet-project-1/services/api/internal/app"
	dbhttp "github.com/dodocheck/go-pet-project-1/services/api/internal/clients/db/http"
	"github.com/dodocheck/go-pet-project-1/services/api/internal/transport/http"
)

func main() {
	dbServiceStr := "http://db-service"
	dbClient := dbhttp.NewDBClient(dbServiceStr)
	service := app.NewService(dbClient)
	httpServer := http.NewHttpServer(service)

	if err := httpServer.StartServer(); err != nil {
		log.Fatal("Failed to start http web server:", err)
	}
}
