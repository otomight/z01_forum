package database

import (
	"forum/internal/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Create client for testing
func InsertSampleClient() {
	// Check if the sample client already exists
	var exists	bool
	var	k		config.StructTablesKeys = config.TableKeys

	query := `
		SELECT EXISTS(SELECT 1 FROM `+k.Clients.Table+` WHERE email = ?)
	`
	err := DB.QueryRow(query, "sample@example.com").Scan(&exists)
	if err != nil {
		log.Printf("Error checking if sample client exists: %v", err)
		return
	}

	// Only insert if it does not exist
	if !exists {
		// Hash the password
		password := "securepassword"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return
		}

		// Insert the client with the hashed password
		query := `
			INSERT INTO `+k.Clients.Table+` (last_name, first_name,
							user_name, email, password, user_role)
			VALUES ('Doe', 'John', 'johndoe',
					'sample@example.com', ?, 'administrator')
		`
		_, err = DB.Exec(query, hashedPassword)
		if err != nil {
			log.Printf("Failed to insert sample client: %v", err)
		} else {
			log.Println("Sample client inserted successfully.")
		}
	} else {
		log.Println("Sample client already exists, skipping insertion.")
	}
}

// inserting post for testing
func InsertSamplePost() {
	// Check if any posts already exist
	var	count	int
	var	query	string
	var	k		config.StructTablesKeys = config.TableKeys

	query = `
		SELECT COUNT(*) FROM `+k.Posts.Table+`
	`
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Failed to count posts: %v", err)
		return
	}
	// If no posts exist, insert the sample post
	if count == 0 {
		query = `
			INSERT INTO `+k.Posts.Table+` (author_id, title, category, tags, content,
							creation_date, update_date, is_deleted)
			VALUES (1, 'Sample Post', 'General', 'other',
					'This is a sample post content.',
					CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0)
		`
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("Failed to insert sample post: %v", err)
		}
	}
}
