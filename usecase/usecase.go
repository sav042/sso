package usecase

import (
	"context"
	ssogrpc "github.com/sav042/sso/usecase/clients/sso/grpc"
	"log/slog"
	"os"
	"time"
)

const userID = "e96ef0f5-724f-43c5-9046-f0c79348be79"

func AppExample() {
	// app code
	log := slog.Logger{}

	// auth via sso
	ssoClient, err := ssogrpc.New(
		context.Background(),
		&log,
		"localhost",
		30*time.Second,
		3,
	)
	if err != nil {
		log.Error("failed to init sso client", err.Error())
		os.Exit(1)
	}

	isAdmin, err := ssoClient.IsAdmin(context.Background(), userID)
	if err != nil {
		panic(err)
	}

	// using auth result
	_ = isAdmin
}
