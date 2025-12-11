package postgres

import (
	"database/sql"
	"log"
	"pet1/models"
)

func initDB() *sql.DB {
	connStr := "postgres://my_user:my_password@localhost:5432/my_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
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
	tasks := []models.TaskImportData{
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
