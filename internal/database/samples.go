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
	var	c		config.ClientsTableKeys

	c = config.TableKeys.Clients
	query := `
		SELECT EXISTS(SELECT 1 FROM `+c.Clients+` WHERE `+c.Email+` = ?);
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
			INSERT INTO `+c.Clients+` (
				`+c.LastName+`, `+c.FirstName+`, `+c.UserName+`,
				`+c.Email+`, `+c.Password+`, `+c.UserRole+`
			)
			VALUES (
				'Doe', 'John', 'johndoe',
				'sample@example.com', ?, 'administrator'
			);
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
	var	p		config.PostsTableKeys

	p = config.TableKeys.Posts
	query = `
		SELECT COUNT(*) FROM `+p.Posts+`;
	`
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Failed to count posts: %v", err)
		return
	}
	// If no posts exist, insert the sample post
	if count == 0 {
		query = `
			INSERT INTO `+p.Posts+` (
				`+p.AuthorID+`, `+p.Title+`, `+p.Category+`,
				`+p.Tags+`, `+p.Content+`, `+p.CreationDate+`,
				`+p.UpdateDate+`, `+p.IsDeleted+`
			)
			VALUES (
				1, 'Sample Post', 'General',
				'other', 'This is a sample post content.', CURRENT_TIMESTAMP,
				CURRENT_TIMESTAMP, 0
			);
		`
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("Failed to insert sample post: %v", err)
		}
	}
}
