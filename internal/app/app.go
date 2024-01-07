package app

import (
	grpcapp "github.com/sav042/sso/internal/app/grpc"
	"github.com/sav042/sso/internal/services/auth"
	"github.com/sav042/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// init storage
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// init auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	// init grpc app
	grpcApp := grpcapp.New(log, grpcPort, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
