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

func (s *Storage) Get(ctx context.Context, projectID uint64) (project models.Project, err error) {
	const op = "storage.postgres.Project"

	query := `SELECT id, name, description FROM project WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, projectID)

	err = row.Scan(&project.Id, &project.Name, &project.Description)

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

func (s *Storage) Update(ctx context.Context, projectID uint64, name, description string) (bool, error) {
	const op = "storage.postgres.Update"

	query := `UPDATE project SET name = $1, description = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, name, description, projectID)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}

	return true, nil
}

func (s *Storage) Delete(ctx context.Context, projectID uint64) (bool, error) {
	const op = "storage.postgres.Delete"

	query := `DELETE FROM project WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, projectID)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil
}

func (s *Storage) AddUser(ctx context.Context, projectId, userId uint64) (bool, error) {
	const op = "storage.postgres.AddUser"

	query := `INSERT INTO users (project_id, user_id) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, projectId, userId)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil

}

func (s *Storage) DeleteUser(ctx context.Context, projectId, userId uint64) (bool, error) {
	const op = "storage.postgres.RemoveUser"

	query := `DELETE FROM users WHERE project_id = $1 AND user_id = $2`

	_, err := s.db.ExecContext(ctx, query, projectId, userId)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil
}

func (s *Storage) Users(ctx context.Context, projectId uint64) (users []models.User, err error) {
	const op = "storage.postgres.Users"

	query := `SELECT project_id, user_id FROM users WHERE project_id = $1`

	rows, err := s.db.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ProjectId, &user.UserId)
		if err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Storage) CreateTask(ctx context.Context, projectId, taskId uint64) (bool, error) {
	const op = "storage.postgres.CreateTask"

	query := `INSERT INTO task (project_id, task_id) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, projectId, taskId)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil
}

func (s *Storage) DeleteTask(ctx context.Context, taskId, projectID uint64) (bool, error) {
	const op = "storage.postgres.DeleteTask"

	query := `DELETE FROM task WHERE task_id = $2 AND project_id = $1`

	_, err := s.db.ExecContext(ctx, query, taskId, projectID)
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}
	return true, nil
}

func (s *Storage) GetTasks(ctx context.Context, projectId uint64) (tasks []models.Task, err error) {
	const op = "storage.postgres.Tasks"

	query := `SELECT project_id, task_id FROM task WHERE project_id = $1`

	rows, err := s.db.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ProjectId, &task.TaskId)
		if err != nil {
			return nil, fmt.Errorf("%s %w", op, err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
