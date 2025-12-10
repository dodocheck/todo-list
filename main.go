package main

import (
	"pet1/services/db/postgres"

	"github.com/k0kubun/pp/v3"
	_ "github.com/lib/pq"
)

func main() {

	db := postgres.Init()

	tasks := postgres.ListAllTasks(db)
	pp.Println(tasks)

	newTaskId := postgres.InsertTask(db, postgres.TaskInsertData{Title: "Тестовый Title", Text: "Тестовый Text"})
	newTask := postgres.GetTask(db, newTaskId)
	pp.Println(newTask)

	postgres.Finish(db)
}
