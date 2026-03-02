package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-refresh-practice/go-refresh-course/cmd/api"
	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/db"
	"github.com/go-refresh-practice/go-refresh-course/service/seed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
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
			config.Envs.DBAddress, // already host:port
			config.Envs.DBName,
		)
	}

	log.Println("Connecting to database...")

	dbConn, err := db.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(dbConn)
	runMigrations(dbConn)

	seed.SeedAdmin(dbConn)

	port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}

addr := "0.0.0.0:" + port

log.Printf("Server running on %s", addr)

server := api.NewAPIServer(addr, dbConn)
if err := server.Run(); err != nil {
	log.Fatal(err)
}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DB: Successfully connected !!")
}


func runMigrations(dbConn *sql.DB) {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

m, err := migrate.NewWithDatabaseInstance(
	"file://migrations",
	"postgres",
	driver,
)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("✅ Database migrations up to date")
}