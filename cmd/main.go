package main

import (
	"log"

	"github.com/go-refresh-practice/go-refresh-course/cmd/api"
	"github.com/go-refresh-practice/go-refresh-course/db"
	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
		
		
	
	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err !=nil {
		log.Fatal(err);
	}
}