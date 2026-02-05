package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/db"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// PRIORITY 1: Check for DATABASE_URL (Fly.io sets this automatically)
	dsn := os.Getenv("DATABASE_URL")
	
	// PRIORITY 2: Fallback to building DSN from config (local development)
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			config.Envs.DBUser,
			config.Envs.DBPassword,
			config.Envs.DBAddress,
			config.Envs.DBName,
		)
	}

	log.Println("Connecting to database for migrations...")

	// Connect using existing helper
	dbConn, err := db.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Migration driver for Postgres
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// "up" or "down" command
	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("✅ Migrations applied successfully!")
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("✅ Migrations rolled back successfully!")
	}
}
