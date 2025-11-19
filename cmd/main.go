package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-refresh-practice/go-refresh-course/cmd/api"
	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/db"
	"github.com/go-refresh-practice/go-refresh-course/service/seed"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Build PostgreSQL DSN from your config
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress, // already host:port
		config.Envs.DBName,
	)

	dbConn, err := db.NewPostgresStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(dbConn)

	seed.SeedAdmin(dbConn)

	server := api.NewAPIServer(":8080", dbConn)
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
