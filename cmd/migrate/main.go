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

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			config.Envs.DBUser,
			config.Envs.DBPassword,
			config.Envs.DBAddress,
			config.Envs.DBName,
		)
	}

	dbConn, err := db.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Specify 'up' or 'down'")
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("✅ Migrations applied successfully!")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		log.Println("✅ Migrations rolled back successfully!")

	default:
		log.Fatal("Unknown command. Use 'up' or 'down'")
	}
}