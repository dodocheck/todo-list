package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type TaskInsertData struct {
	Title string
	Text  string
}

type TaskReceiveData struct {
	Id         int64
	Title      string
	Text       string
	Finished   bool
	CreatedAt  time.Time
	FinishedAt *time.Time
}

func Init() *sql.DB {
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

func Finish(db *sql.DB) {
	db.Close()
}

func createTasksTable(db *sql.DB) {
	dropQuery := `drop table tasks`
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

func InsertTask(db *sql.DB, task TaskInsertData) int64 {
	query := `insert into tasks (title,text) values ($1,$2) returning id`

	var id int64
	if err := db.QueryRow(query, task.Title, task.Text).Scan(&id); err != nil {
		log.Fatal(err)
	}

	return id
}

func GetTask(db *sql.DB, id int64) TaskReceiveData {
	var taskToReturn TaskReceiveData
	if err := db.QueryRow("select id, title, text, finished, created_at, finished_at from tasks where id = $1", id).Scan(
		&taskToReturn.Id,
		&taskToReturn.Title,
		&taskToReturn.Text,
		&taskToReturn.Finished,
		&taskToReturn.CreatedAt,
		&taskToReturn.FinishedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with id = %d", id)
		}
		log.Fatal(err)
	}

	return taskToReturn
}

func ListAllTasks(db *sql.DB) []TaskReceiveData {
	sliceToReturn := make([]TaskReceiveData, 0)

	rows, err := db.Query("select id, title, text, finished, created_at, finished_at from tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var task TaskReceiveData
		if err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Text,
			&task.Finished,
			&task.CreatedAt,
			&task.FinishedAt); err != nil {
			log.Fatal(err)
		}
		sliceToReturn = append(sliceToReturn, task)
	}

	return sliceToReturn
}
