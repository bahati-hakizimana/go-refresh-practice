package main

import (
	"database/sql"
	"log"

	"github.com/go-refresh-practice/go-refresh-course/cmd/api"
	"github.com/go-refresh-practice/go-refresh-course/config"
	"github.com/go-refresh-practice/go-refresh-course/db"
	"github.com/go-refresh-practice/go-refresh-course/service/seed"
	"github.com/go-sql-driver/mysql"
)

func main() {

	 dbConn, err := db.NewMySQLStorage(mysql.Config{

		User:        config.Envs.DBUser,
		Passwd:    config.Envs.DBPassword,
		Addr:        config.Envs.DBAddress,
		DBName:       config.Envs.DBName,
		Net:         "tcp",
		AllowNativePasswords: true,
		ParseTime:   true,
	})
		

	if err != nil {
		log.Fatal(err)
	}
		initStorage( dbConn)
		seed.SeedAdmin( dbConn)
	
	server := api.NewAPIServer(":8080",  dbConn)
	if err := server.Run(); err !=nil {
		log.Fatal(err);
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DB: Successfully connected !!")
}