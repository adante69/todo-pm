package ProjectManager

import (
	"context"
	projv1 "github.com/adante69/todo-protos/gen/go/project"
	"google.golang.org/grpc"
	"todo-project/internal/domain/models"
)

type ProjectManager interface {
	Get(ctx context.Context, id uint64) (models.Project, error)
	Create(ctx context.Context, name, description string) (uint64, error)
	Delete(ctx context.Context, id uint64) (bool, error)
	Update(ctx context.Context, id uint64, name, description string) (bool, error)
}

type TaskManager interface {
	CreateTask(ctx context.Context, projectId, taskId uint64) (bool, error)
	DeleteTask(ctx context.Context, projectId, taskId uint64) (bool, error)
	GetTasks(ctx context.Context, projectId uint64) ([]models.Task, error)
}

type UserManager interface {
	AddUser(ctx context.Context, projectId, userId uint64) (bool, error)
	DeleteUser(ctx context.Context, projectId, userId uint64) (bool, error)
	GetUsers(ctx context.Context, projectId uint64) ([]models.User, error)
}

type serverAPI struct {
	projv1.UnimplementedProjectServer
	ProjectManager ProjectManager
	TaskManager    TaskManager
	UserManager    UserManager
}

func Register(gRpc *grpc.Server, taskManager TaskManager,
	projectManager ProjectManager, userManager UserManager) {
	projv1.RegisterProjectServer(gRpc, &serverAPI{
		ProjectManager: projectManager,
		TaskManager:    taskManager,
		UserManager:    userManager})
}
func (s *serverAPI) Get(ctx context.Context, req *projv1.GetRequest) (*projv1.GetResponse, error) {
	project, err := s.ProjectManager.Get(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &projv1.GetResponse{
		ProjectId:   project.Id,
		Name:        project.Name,
		Description: project.Description,
	}, nil
}

func (s *serverAPI) Create(ctx context.Context, req *projv1.CreateRequest) (*projv1.CreateResponse, error) {
	id, err := s.ProjectManager.Create(ctx, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &projv1.CreateResponse{ProjectId: id}, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *projv1.DeleteRequest) (*projv1.DeleteResponse, error) {
	_, err := s.ProjectManager.Delete(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &projv1.DeleteResponse{Check: true}, nil
}

func (s *serverAPI) Update(ctx context.Context, req *projv1.UpdateRequest) (*projv1.UpdateResponse, error) {
	_, err := s.ProjectManager.Update(ctx, req.ProjectId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &projv1.UpdateResponse{Check: true}, nil

}

func (s *serverAPI) AddTask(ctx context.Context, req *projv1.AddTaskRequest) (*projv1.AddTaskResponse, error) {
	msg, err := s.TaskManager.CreateTask(ctx, req.ProjectId, req.TaskId)
	if err != nil {
		return nil, err
	}
	return &projv1.AddTaskResponse{Check: msg}, nil
}

func (s *serverAPI) DeleteTask(ctx context.Context, req *projv1.DeleteTaskRequest) (*projv1.DeleteTaskResponse, error) {
	msg, err := s.TaskManager.DeleteTask(ctx, req.ProjectId, req.TaskId)
	if err != nil {
		return nil, err
	}
	return &projv1.DeleteTaskResponse{Check: msg}, nil
}

func (s *serverAPI) AddUser(ctx context.Context, req *projv1.AddUserRequest) (*projv1.AddUserResponse, error) {
	msg, err := s.UserManager.AddUser(ctx, req.ProjectId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &projv1.AddUserResponse{Check: msg}, nil
}

func (s *serverAPI) DeleteUser(ctx context.Context, req *projv1.DeleteUserRequest) (*projv1.DeleteUserResponse, error) {
	msg, err := s.UserManager.DeleteUser(ctx, req.ProjectId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &projv1.DeleteUserResponse{Check: msg}, nil
}

func (s *serverAPI) GetUsers(ctx context.Context, req *projv1.GetUsersRequest) (*projv1.GetUsersResponse, error) {
	use, err := s.UserManager.GetUsers(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	var users []uint64
	for _, u := range use {
		users = append(users, u.UserId)
	}
	return &projv1.GetUsersResponse{UserId: users}, nil
}

func (s *serverAPI) GetTasks(ctx context.Context, req *projv1.GetTasksRequest) (*projv1.GetTasksResponse, error) {
	task, err := s.TaskManager.GetTasks(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	var tasks []uint64
	for _, t := range task {
		tasks = append(tasks, t.TaskId)
	}
	return &projv1.GetTasksResponse{
		TaskId: tasks,
	}, nil
}
