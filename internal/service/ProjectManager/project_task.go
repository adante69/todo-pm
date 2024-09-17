package ProjectManager

import (
	"context"
	"log/slog"
	"todo-project/internal/domain/models"
)

type TaskManager struct {
	log          *slog.Logger
	TaskSaver    TaskSaver
	TaskProvider TaskProvider
	TaskRemover  TaskRemover
}

type TaskSaver interface {
	CreateTask(ctx context.Context, projectId, taskId uint64) (bool, error)
}

type TaskProvider interface {
	GetTasks(ctx context.Context, projectID uint64) ([]models.Task, error)
}

type TaskRemover interface {
	DeleteTask(ctx context.Context, taskID, projectID uint64) (bool, error)
}

func NewTaskManager(
	log *slog.Logger,
	TaskSaver TaskSaver,
	TaskProvider TaskProvider,
	TaskRemover TaskRemover,
) *TaskManager {
	return &TaskManager{
		log:          log,
		TaskSaver:    TaskSaver,
		TaskProvider: TaskProvider,
		TaskRemover:  TaskRemover,
	}
}

// CreateTask creates new task
func (m *TaskManager) CreateTask(ctx context.Context, projectId, taskId uint64) (bool, error) {
	id, err := m.TaskSaver.CreateTask(ctx, projectId, taskId)
	if err != nil {
		return id, err
	}
	return id, nil
}

// RemoveTask removes task
func (m *TaskManager) DeleteTask(ctx context.Context, taskID, projectID uint64) (bool, error) {
	check, err := m.TaskRemover.DeleteTask(ctx, taskID, projectID)
	if err != nil {
		return check, err
	}
	return check, nil
}

// Tasks returns tasks
func (m *TaskManager) GetTasks(ctx context.Context, projectID uint64) ([]models.Task, error) {
	tasks, err := m.TaskProvider.GetTasks(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
