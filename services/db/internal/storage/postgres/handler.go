package postgres

import (
	"database/sql"
	"log"

	"github.com/dodocheck/go-pet-project-1/pb"
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

func (pc *PostgresController) AddTask(task pb.TaskImportData) (pb.TaskExportData, error) {
	query := `insert into tasks (title,text) values ($1,$2) returning id, title, text, finished, created_at, finished_at`

	var createdTask pb.TaskExportData
	if err := pc.db.QueryRow(query, task.Title, task.Text).Scan(
		&createdTask.Id,
		&createdTask.Title,
		&createdTask.Text,
		&createdTask.Finished,
		&createdTask.CreatedAt,
		&createdTask.FinishedAt); err != nil {
		return pb.TaskExportData{}, err
	}

	return createdTask, nil
}

func (pc *PostgresController) DeleteTask(id int) error {
	if _, err := pc.db.Exec("delete from tasks where id = $1", id); err != nil {
		return err
	}

	return nil
}

func (pc *PostgresController) ListAllTasks() ([]pb.TaskExportData, error) {
	sliceToReturn := make([]pb.TaskExportData, 0)

	rows, err := pc.db.Query("select id, title, text, finished, created_at, finished_at from tasks order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task pb.TaskExportData
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

func (pc *PostgresController) MarkTaskFinished(id int) (pb.TaskExportData, error) {
	query := `update tasks 
        set finished = true, 
        finished_at = NOW() 
        where id = $1 
        returning id, title, text, finished, created_at, finished_at`

	var updatedTask pb.TaskExportData

	if err := pc.db.QueryRow(query, id).Scan(
		&updatedTask.Id,
		&updatedTask.Title,
		&updatedTask.Text,
		&updatedTask.Finished,
		&updatedTask.CreatedAt,
		&updatedTask.FinishedAt); err != nil {
		return pb.TaskExportData{}, err
	}

	return updatedTask, nil
}
