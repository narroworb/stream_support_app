package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"stream_app/internal/migrate"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var PG *sql.DB

func InitPostgres() {
	dsn := os.Getenv("POSTGRES_DSN")

	var err error
	PG, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	if err := PG.Ping(); err != nil {
		log.Fatalf("Postgres ping error: %v", err)
	}

	fmt.Println("Connected to Postgres")

	migrate.RunPostgresMigrations(PG)

	fmt.Println("Migrations applied")
}
