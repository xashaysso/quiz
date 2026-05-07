package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func execSQL(ctx context.Context, pool *pgxpool.Pool, sql string){
	_, err := pool.Exec(ctx, sql);
	if err != nil{
		log.Fatalf("Error creating table: %v\nSQL: %s", err, sql);
	}
}

func CreateTables(ctx context.Context, pool *pgxpool.Pool) {
	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS quiz(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	)`);

	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS questions(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		quiz_id INT NOT NULL REFERENCES quiz(id) ON DELETE CASCADE
	)`);

	execSQL(ctx, pool, `CREATE TABLE IF NOT EXISTS answers(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		correct BOOLEAN DEFAULT false
	)`);

	log.Println("All tables created succesfully.");
}

func Serve() *pgxpool.Pool {
	ctx := context.Background();

	DB_URL := os.Getenv("DB_URL");

	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil{
		log.Fatalf("Unable to connect to db: %v", err);
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err);
	}

	log.Printf("Connected to db succesfully");

	CreateTables(ctx, pool);

	return pool;
}