package seed

import (
	"database/sql"
	"log"

	"github.com/go-refresh-practice/go-refresh-course/service/auth"
)

func SeedAdmin(db *sql.DB){

	// check if admin exists
    var count int
	 err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
    if err != nil {
        log.Println("Error checking admin:", err)
        return
    }

	  if count > 0 {
        log.Println("Admin already seeded")
        return
    }

	password, _ := auth.HashPassword("Rwiyereka123")
	  _, err = db.Exec(`
        INSERT INTO users (firstName, lastName, email, password, role)
        VALUES (?, ?, ?, ?, ?)`,
        "System", "Admin", "clairesblapt@gmail.com ", password, "admin",
    )

    if err != nil {
        log.Println("Failed to seed admin:", err)
        return
    }

	log.Println("Admin user seeded successfully")
     	
}