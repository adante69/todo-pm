package ProjectManager

import (
	"context"
	"log/slog"
	"todo-project/internal/domain/models"
)

type UsersManager struct {
	log          *slog.Logger
	UserAdder    UserAdder
	UserProvider UserProvider
	UserRemover  UserRemover
}
type UserAdder interface {
	AddUser(ctx context.Context, projectId, userId uint64) (bool, error)
}
type UserProvider interface {
	Users(ctx context.Context, projectId uint64) ([]models.User, error)
}

type UserRemover interface {
	DeleteUser(ctx context.Context, projectId, userId uint64) (bool, error)
}

func NewUsersManager(
	log *slog.Logger,
	UserAdder UserAdder,
	UserProvider UserProvider,
	UserRemover UserRemover,
) *UsersManager {
	return &UsersManager{
		log:          log,
		UserAdder:    UserAdder,
		UserProvider: UserProvider,
		UserRemover:  UserRemover,
	}
}

// AddUser adds user to project
func (m *UsersManager) AddUser(ctx context.Context, projectId, userId uint64) (bool, error) {
	check, err := m.UserAdder.AddUser(ctx, projectId, userId)
	if err != nil {
		return check, err
	}
	return check, nil
}

// RemoveUser removes user from project
func (m *UsersManager) DeleteUser(ctx context.Context, projectId, userId uint64) (bool, error) {
	check, err := m.UserRemover.DeleteUser(ctx, projectId, userId)
	if err != nil {
		return check, err
	}
	return check, nil
}

// Users returns all users of project
func (m *UsersManager) GetUsers(ctx context.Context, projectId uint64) ([]models.User, error) {
	users, err := m.UserProvider.Users(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return users, nil
}
