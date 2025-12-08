package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func execSQL(ctx context.Context, conn *pgx.Conn, sql string){
	_, err := conn.Exec(ctx, sql);
	if err != nil{
		log.Fatalf("Error creating table: %v\nSQL: %s", err, sql);
	}
}

func CreateTables(ctx context.Context, conn *pgx.Conn) {
	execSQL(ctx, conn, `CREATE TABLE IF NOT EXISTS quiz(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT
	)`);

	execSQL(ctx, conn, `CREATE TABLE IF NOT EXISTS questions(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		quiz_id INT NOT NULL REFERENCES quiz(id) ON DELETE CASCADE
	)`);

	execSQL(ctx, conn, `CREATE TABLE IF NOT EXISTS answers(
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		correct BOOLEAN DEFAULT false
	)`);

	log.Println("All tables created succesfully.");
}

func Serve() *pgx.Conn {
	ctx := context.Background();

	DB_URL := os.Getenv("DB_URL");

	conn, err := pgx.Connect(ctx, DB_URL);
	if err != nil{
		log.Fatalf("Unable to connect to db: %v", err);
	}
	log.Printf("Connected to db succesfully");

	CreateTables(ctx, conn);

	return conn;
}