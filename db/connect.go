package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func execSQL(ctx context.Context, pool *pgxpool.Pool, sql string) {
	_, err := pool.Exec(ctx, sql)
	if err != nil {
		slog.Error("failed to execute SQL query",
			slog.Any("err", err),
			slog.String("sql", sql),
		)
		os.Exit(1)
	}
}

func CreateTables(ctx context.Context, pool *pgxpool.Pool) {
	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	)`)

	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS quiz(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		creator_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
	)`)

	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS questions(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		quiz_id INT NOT NULL REFERENCES quiz(id) ON DELETE CASCADE
	)`)

	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS answers(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		correct BOOLEAN DEFAULT false
	)`)

	slog.Info("all tables created succesfully.")
}

func Serve() *pgxpool.Pool {
	ctx := context.Background()

	DB_URL := os.Getenv("DB_URL")

	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil {
		slog.Error("unable to create database pool", slog.Any("err", err))
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", slog.Any("err", err))
	}

	slog.Info("connected to db successfully")

	CreateTables(ctx, pool)

	return pool
}

func NewRedisClient(addr string) *redis.Client {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("redis connection failed", slog.Any("err", err), slog.String("addr", addr))
		os.Exit(1)
	}
	slog.Info("connected to redis successfully", slog.String("addr", addr))

	return rdb
}
