package seed

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-refresh-practice/go-refresh-course/service/auth"
)

func SeedAdmin(db *sql.DB) {
	// Check if admin exists
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

	// Get admin credentials from environment variables
	firstName := os.Getenv("ADMIN_FIRST_NAME")
	lastName := os.Getenv("ADMIN_LAST_NAME")
	email := os.Getenv("ADMIN_EMAIL")
	plainPassword := os.Getenv("ADMIN_PASSWORD")

	// Validate environment variables
	if firstName == "" || lastName == "" || email == "" || plainPassword == "" {
		log.Println("Missing admin credentials in environment variables")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(plainPassword)
	if err != nil {
		log.Println("Error hashing admin password:", err)
		return
	}

	// Insert admin user
	_, err = db.Exec(`
		INSERT INTO users (first_name, last_name, email, password, role)
		VALUES ($1, $2, $3, $4, $5)`,
		firstName, lastName, email, hashedPassword, "admin",
	)

	if err != nil {
		log.Println("Failed to seed admin:", err)
		return
	}

	log.Println("Admin user seeded successfully")
}
