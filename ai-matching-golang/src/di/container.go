package di

import (
	db "ai-matching/db/sqlc"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Container struct {
	DB      *sqlx.DB
	Queries db.Querier
}

func NewContainer() *Container {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		sslMode := os.Getenv("DB_SSL_MODE")

		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}
		if sslMode == "" {
			sslMode = "disable"
		}

		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	queries := db.New(sqlDB)

	return &Container{
		DB:      sqlxDB,
		Queries: queries,
	}
}
