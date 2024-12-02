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
		_, err = insertInto(InsertIntoQuery{
			Table:	c.Clients,
			Keys: []string{
				c.LastName, c.FirstName, c.UserName,
				c.Email, c.Password, c.UserRole,
			},
			Values: [][]any{{
				"Doe", "John", "johndoe",
				"sample@example.com", hashedPassword, "administrator",
			}},
		})

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
	var	p		config.PostsTableKeys

	p = config.TableKeys.Posts
	count = countRows(p.Posts)
	// If no posts exist, insert the sample post
	if count == 0 {
		_, err := insertInto(InsertIntoQuery{
			Table: p.Posts,
			Keys: []string{p.AuthorID, p.Title, p.Content, p.IsDeleted},
			Values: [][]any{{
				1, "Sample Post", "This is a sample post content.", 0,
			}},
		})
		if err != nil {
			log.Printf("Failed to insert sample post: %v", err)
		}
	}
}
