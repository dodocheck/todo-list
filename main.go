package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

type Task struct {
	Id        int64
	Title     string
	Text      string
	CreatedAt time.Time
}

func main() {
	dsn := "postgres://my_user:my_password@localhost:5432/my_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Error while opening my_db", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("Ping failed", err)
	}

	fmt.Println("Ping success!")

	createTableSQL := `
create table if not exists tasks (
	id bigserial primary key,
	title text not null,
	text text,
	created_at timestamp not null default NOW()
);`

	if _, err := db.Exec(createTableSQL); err != nil {
		fmt.Println("create table error:", err)
		return
	}

	fmt.Println("table tasks created")

	var insertedId int
	insertSQL := `insert into tasks (title,text) values ($1,$2) returning id`
	if err := db.QueryRow(insertSQL, "Do hw", "do math hw till tomorrow").Scan(&insertedId); err != nil {
		fmt.Println("error while inserting value:", err)
		return
	}
	fmt.Println("Successfully inserted row with id:", insertedId)

	rows, err := db.Query("select * from tasks")
	if err != nil {
		fmt.Println("Error while getting all rows:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Text, &task.CreatedAt); err != nil {
			fmt.Println("Failed to read row:", err)
		}
		pp.Println(task)
	}
}
