package postgres

import "database/sql"

type PostgresController struct {
	db *sql.DB
}

func NewPostgresController() *PostgresController {
	return &PostgresController{db: initDB()}
}
