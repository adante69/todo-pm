package ProjectManager

import (
	"context"
	"log/slog"
	"todo-project/internal/domain/models"
)

type ProjectManager struct {
	log             *slog.Logger
	projectSaver    ProjectSaver
	projectProvider ProjectProvider
	projectRemover  ProjectRemover
}

type ProjectSaver interface {
	Create(ctx context.Context, name, description string) (uint64, error)
	Update(ctx context.Context, projectId uint64, name, description string) (bool, error)
}

type ProjectProvider interface {
	Get(ctx context.Context, taskID uint64) (models.Project, error)
}

type ProjectRemover interface {
	Delete(ctx context.Context, taskID uint64) (bool, error)
}

func NewProjectManager(
	log *slog.Logger,
	projectSaver ProjectSaver,
	projectProvider ProjectProvider,
	projectRemover ProjectRemover,
) *ProjectManager {
	return &ProjectManager{
		log:             log,
		projectSaver:    projectSaver,
		projectProvider: projectProvider,
		projectRemover:  projectRemover,
	}
}

// Create creates new project
func (m *ProjectManager) Create(ctx context.Context, name, description string) (uint64, error) {
	id, err := m.projectSaver.Create(ctx, name, description)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update updates project
func (m *ProjectManager) Update(ctx context.Context, projectId uint64, name, description string) (bool, error) {
	check, err := m.projectSaver.Update(ctx, projectId, name, description)
	if err != nil {
		return check, err
	}
	return check, nil
}

// Delete deletes project
func (m *ProjectManager) Delete(ctx context.Context, id uint64) (bool, error) {
	check, err := m.projectRemover.Delete(ctx, id)
	if err != nil {
		return check, err
	}
	return check, nil
}

// Get returns project by id
func (m *ProjectManager) Get(ctx context.Context, id uint64) (models.Project, error) {
	project, err := m.projectProvider.Get(ctx, id)
	if err != nil {
		return models.Project{}, err
	}
	return project, nil
}
