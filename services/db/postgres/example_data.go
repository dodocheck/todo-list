package postgres

import (
	"database/sql"
	"log"
)

func seedTasks(db *sql.DB) {
	tasks := []TaskInsertData{
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
