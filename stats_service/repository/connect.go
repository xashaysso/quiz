package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dbURL string) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		slog.Error("failed to open db connection for migrations", slog.Any("err", err))
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("failed to set goose dialect", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("migration up failed", slog.Any("err", err))
		os.Exit(1)
	}
	slog.Info("database migrations applied successfully")
}

func Serve() *pgxpool.Pool {
	ctx := context.Background()

	DB_URL := os.Getenv("DB_URL")

	RunMigrations(DB_URL)

	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil {
		slog.Error("unable to create database pool", slog.Any("err", err))
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", slog.Any("err", err))
		os.Exit(1)
	}

	slog.Info("connected to db successfully")

	return pool
}
