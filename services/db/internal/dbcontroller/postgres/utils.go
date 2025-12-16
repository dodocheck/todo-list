package postgres

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/dodocheck/go-pet-project-1/services/db/pb"
	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	pUser := os.Getenv("POSTGRES_USER")
	pPassword := os.Getenv("POSTGRES_PASSWORD")
	pDb := os.Getenv("POSTGRES_DB")
	pPort := os.Getenv("POSTGRES_PORT")
	connStr := "postgres://" + pUser + ":" + pPassword + "@postgres:" + pPort + "/" + pDb + "?sslmode=disable"

	var db *sql.DB
	var err error

	connRetries := 15
	for range connRetries {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			log.Println("Connected to Postgres!")
			break
		}
		log.Println("Postgres not ready yet, retrying...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to Postgres:", err)
	}

	createTasksTable(db)

	seedTasks(db)

	return db
}

func createTasksTable(db *sql.DB) {
	dropQuery := `drop table if exists tasks`
	if _, err := db.Exec(dropQuery); err != nil {
		log.Fatal(err)
		return
	}

	createQuery := `create table if not exists tasks (
                id bigserial primary key,
                title varchar(20) not null,
                text varchar(100),
                finished bool default false,
                created_at timestamp not null default NOW(),
                finished_at timestamp default NULL);`

	if _, err := db.Exec(createQuery); err != nil {
		log.Fatal(err)
		return
	}
}

func seedTasks(db *sql.DB) {
	tasks := []pb.TaskImportData{
		{Title: "Помыть посуду", Text: "После ужина на кухне"},
		{Title: "Сходить в зал", Text: "Тренировка спины и ног"},
		{Title: "Позвонить маме", Text: "Уточнить планы на выходные"},
		{Title: "Почитать Go", Text: "1 глава про горутины"},
		{Title: "Купить продукты", Text: "Молоко, хлеб, сыр, яйца"},
		{Title: "Оплатить счёт", Text: "Коммуналка до пятницы"},
		{Title: "Написать таск", Text: "Протестировать InsertTask и GetTask"},
		{Title: "Убраться дома", Text: "Пылесос и влажная уборка"},
		{Title: "Погулять", Text: "30 минут прогулки без телефона"},
		{Title: "Сериал", Text: "Посмотреть 1 серию вечером"},
	}

	for _, t := range tasks {
		_, err := db.Exec(
			`insert into tasks (title, text) values ($1, $2)`,
			t.Title, t.Text,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}
