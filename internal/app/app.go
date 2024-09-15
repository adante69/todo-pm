package app

import (
	"log/slog"
	grpcapp "todo-project/internal/app/grpc"
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
	storage, err := postgres.New(dsn)
	if err != nil {
		panic(err)
	}
}
