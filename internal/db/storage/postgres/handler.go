package postgres

import (
	"database/sql"
	"log"
	"pet1/pkg/contracts"

	_ "github.com/lib/pq"
)

type PostgresController struct {
	db *sql.DB
}

func NewPostgresController() *PostgresController {
	return &PostgresController{db: initDB()}
}

func (pc *PostgresController) Close() {
	pc.db.Close()
}

func (pc *PostgresController) AddTask(task contracts.TaskImportData) (contracts.TaskExportData, error) {
	query := `insert into tasks (title,text) values ($1,$2) returning id, title, text, finished, created_at, finished_at`

	var createdTask contracts.TaskExportData
	if err := pc.db.QueryRow(query, task.Title, task.Text).Scan(
		&createdTask.Id,
		&createdTask.Title,
		&createdTask.Text,
		&createdTask.Finished,
		&createdTask.CreatedAt,
		&createdTask.FinishedAt); err != nil {
		return contracts.TaskExportData{}, err
	}

	return createdTask, nil
}

func (pc *PostgresController) DeleteTask(id int) error {
	if _, err := pc.db.Exec("delete from tasks where id = $1", id); err != nil {
		return err
	}

	return nil
}

func (pc *PostgresController) ListAllTasks() ([]contracts.TaskExportData, error) {
	sliceToReturn := make([]contracts.TaskExportData, 0)

	rows, err := pc.db.Query("select id, title, text, finished, created_at, finished_at from tasks order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task contracts.TaskExportData
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

	return sliceToReturn, nil
}

func (pc *PostgresController) MarkTaskFinished(id int) (contracts.TaskExportData, error) {
	query := `update tasks 
        set finished = true, 
        finished_at = NOW() 
        where id = $1 
        returning id, title, text, finished, created_at, finished_at`

	var updatedTask contracts.TaskExportData

	if err := pc.db.QueryRow(query, id).Scan(
		&updatedTask.Id,
		&updatedTask.Title,
		&updatedTask.Text,
		&updatedTask.Finished,
		&updatedTask.CreatedAt,
		&updatedTask.FinishedAt); err != nil {
		return contracts.TaskExportData{}, err
	}

	return updatedTask, nil
}
