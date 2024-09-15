package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, pm *taskManager.TaskManager, port int) *App {
	gRPCServer := grpc.NewServer()
	tmsgrpc.Register(gRPCServer, pm)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"
	log := a.log.With(slog.String("op", op))

	log.Info("start")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}
	log.Info("Running grpc Server on port %d", a.port)

	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}
	return nil
}

func (a *App) Stop() {
	const op = "app.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping server on port %d", a.port)
	a.gRPCServer.GracefulStop()
}
