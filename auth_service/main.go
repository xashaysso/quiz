package main

import (
	"auth/pkg/authv1"
	"auth/repository"
	"auth/repository/pg"
	"auth/repository/redis"
	"auth/service"
	"auth/transport"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("No .env file found")
	}

	// logger
	var logger *slog.Logger

	if os.Getenv("APP_ENV") == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)

	slog.Info("Logger service has started", slog.String("env", os.Getenv("APP_ENV")))

	PORT := os.Getenv("PORT")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")

	// postgre init
	globalPool := repository.Serve()
	defer globalPool.Close()

	// repositories
	userRepo := pg.NewUserRepo(globalPool)

	// redis init
	rdb := repository.NewRedisClient(REDIS_ADDR)

	sessionRepo := redis.NewSessionRepository(rdb)

	authService := service.NewAuthService(userRepo, sessionRepo)

	// transport
	authH := transport.NewAuthHandler(authService)

	// grpc server
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		slog.Error("failed to listen tcp port", slog.String("port", PORT), slog.Any("err", err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authH)

	go func() {
		slog.Info("gRPC server is running", slog.String("port", PORT))
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server crashed", slog.Any("err", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down server...")
	grpcServer.GracefulStop()
	slog.Info("server exited cleanly")
}
