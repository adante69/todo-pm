package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"todo-project/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(source string) (*Storage, error) {
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) Project(ctx context.Context, projectID uint64) (project models.Project, err error) {
	const op = "storage.postgres.Project"

	query := `SELECT id, name, description FROM project WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, projectID)

	err = row.Scan(&project.name, &project.description, &project.id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, fmt.Errorf("%s %w", op, err)
		}
		return models.Project{}, fmt.Errorf("%s %w", op, err)
	}
	return project, nil

}

func (s *Storage) Create(ctx context.Context, name, description string) (id uint64, err error) {
	const op = "storage.postgres.Create"

	query := `INSERT INTO project (name, description) VALUES ($1, $2) RETURNING id`

	err = s.db.QueryRowContext(ctx, query, name, description).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("%s %w", op, err)
	}
	return id, nil
}
