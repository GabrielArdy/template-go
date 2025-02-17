package main

import (
	"context"
	"errors"
	"fmt"
	"go-scratch/config"
	"go-scratch/generated"
	"go-scratch/internal/handler"
	"go-scratch/internal/repository"
	"go-scratch/internal/services"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	config.Load(ctx)
	e := config.LoadEcho()

	ar := repository.NewAuthRepository(config.Cli.MongoDB, "auth")
	ur := repository.NewUserRepository(config.Cli.MongoDB, "users")
	lr := repository.NewActivityRepository(config.Cli.MongoDB, "activity_logs")
	ls := services.NewLoggingService(lr)
	atr := repository.NewAttendanceRepository(config.Cli.MongoDB, "attendance_logs")
	ats := services.NewAttendanceService(atr, ls, config.Cli.Redis)
	uas := services.NewUserAuthService(ur, ar, ls, config.Cli.Redis)
	var server generated.ServerInterface = handler.NewHandler(uas, ats)
	generated.RegisterHandlers(e, server)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", config.Conf.Server.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("error when starting up server", slog.Any("err", err.Error()))
			os.Exit(-1)
		}
	}()

	wait := config.GracefulShutdown(context.Background(), 2*time.Second, map[string]func(ctx context.Context) error{
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	<-wait

	slog.Info("gracefully shutdown server")
}
