package main

import (
	"log"
	"pet1/services/api"
	"pet1/services/db/postgres"
	web "pet1/services/web/http"

	_ "github.com/lib/pq"
)

func main() {
	postgresController := postgres.NewPostgresController()
	toDoList := api.NewToDoList(postgresController)
	httpServer := web.NewHttpServer(toDoList)

	if err := httpServer.StartServer(); err != nil {
		log.Fatal("Failed to start http web server:", err)
	}
}
