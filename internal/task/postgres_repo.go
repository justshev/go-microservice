package task

import (
	"context"
	"database/sql"
)


type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) List(ctx context.Context) ([]Task, error) {
	rows,err := r.db.QueryContext(ctx,` 
		SELECT id, name, completed, created_at
		FROM tasks
		ORDER BY id
	`)
	if err != nil { 
		return nil, err
	}

	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
} 

func (r *PostgresRepo) Create(ctx context.Context, name string) (Task, error) {
	var t Task
	err := r.db.QueryRowContext(ctx,`
		INSERT INTO tasks (name, completed, created_at)
		VALUES ($1, false, NOW())
		RETURNING id, name, completed, created_at
	`, name).Scan(&t.ID, &t.Name, &t.Completed, &t.CreatedAt)
	if err != nil {
		return Task{}, err
	}

	return t, err
}