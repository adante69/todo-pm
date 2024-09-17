package app

import (
	"log/slog"
	grpcapp "todo-project/internal/app/grpc"
	"todo-project/internal/service/ProjectManager"
	"todo-project/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	dsn string,
) *App {
	storage, err := postgres.NewStorage(dsn)
	if err != nil {
		panic(err)
	}

	pm := ProjectManager.NewProjectManager(log, storage, storage, storage)
	tm := ProjectManager.NewTaskManager(log, storage, storage, storage)
	um := ProjectManager.NewUsersManager(log, storage, storage, storage)

	grpcServer := grpcapp.New(log, tm, pm, um, grpcPort)

	return &App{
		GRPCServer: grpcServer,
	}
}
