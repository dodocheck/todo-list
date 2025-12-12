package main

import (
	"log"
	"pet1/internal/api/app"
	dbhttp "pet1/internal/api/clients/db/http"
	"pet1/internal/api/transport/http"
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
