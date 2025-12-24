package postgres

import (
	"context"

	"github.com/dodocheck/go-pet-project-1/services/db/internal/models"
)

func (pc *PostgresController) Close() error {
	return pc.db.Close()
}

func (pc *PostgresController) AddTask(ctx context.Context, task models.TaskImportData) (models.TaskExportData, error) {
	query := `insert into tasks (title,text) values ($1,$2) returning id, title, text, finished, created_at, finished_at`

	var createdTask models.TaskExportData
	if err := pc.db.QueryRowContext(ctx, query, task.Title, task.Text).Scan(
		&createdTask.Id,
		&createdTask.Title,
		&createdTask.Text,
		&createdTask.Finished,
		&createdTask.CreatedAt,
		&createdTask.FinishedAt); err != nil {
		return models.TaskExportData{}, err
	}

	return createdTask, nil
}

func (pc *PostgresController) DeleteTask(ctx context.Context, id int) error {
	_, err := pc.db.ExecContext(ctx, "delete from tasks where id = $1", id)

	return err
}

func (pc *PostgresController) ListAllTasks(ctx context.Context) ([]models.TaskExportData, error) {
	sliceToReturn := make([]models.TaskExportData, 0)

	rows, err := pc.db.QueryContext(ctx, "select id, title, text, finished, created_at, finished_at from tasks order by id")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var task models.TaskExportData
		if err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Text,
			&task.Finished,
			&task.CreatedAt,
			&task.FinishedAt); err != nil {
			return nil, err
		}
		sliceToReturn = append(sliceToReturn, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sliceToReturn, nil
}

func (pc *PostgresController) MarkTaskFinished(ctx context.Context, id int) (models.TaskExportData, error) {
	query := `update tasks 
        set finished = true, 
        finished_at = NOW() 
        where id = $1 
        returning id, title, text, finished, created_at, finished_at`

	var updatedTask models.TaskExportData

	if err := pc.db.QueryRowContext(ctx, query, id).Scan(
		&updatedTask.Id,
		&updatedTask.Title,
		&updatedTask.Text,
		&updatedTask.Finished,
		&updatedTask.CreatedAt,
		&updatedTask.FinishedAt); err != nil {
		return models.TaskExportData{}, err
	}

	return updatedTask, nil
}
