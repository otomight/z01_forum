package database

import (
	"database/sql"
	"log"
)

// Client CRUD operations
func CreateClient(client *Client) error {
	_, err := DB.Exec(`
		INSERT INTO Clients (last_name, first_name, user_name, email, password, avatar, bith_date)
		VALUES 	(?, ?, ?, ?, ?, ?, ?)
	`, client.LastName, client.FirstName, client.UserName, client.Email, client.Password, client.Avatar, client.BirthDate)
	return err
}

func GetClientByID(userID int) (*Client, error) {
	row := DB.QueryRow(`
		SELECT user_id, last_name, first_name, user_name, email, avatar, birth_date, creation_date, update_date, deletion_date
		FROM Clients WERE user_id = ?
	`, userID)
	var client Client
	err := row.Scan(&client.UserID, &client.LastName, &client.FirstName, &client.UserName, &client.Email,
		&client.BirthDate, &client.CreationDate, &client.UpdateDate, &client.DeletionDate)
	return &client, err
}

// Retrieve client from database by their Email
func GetClientByUsernameOrEmail(email string) (*Client, error) {
	var client Client

	//Query to find user by Email
	query := `SELECT user_id, last_name, first_name, email, password, avatar, birthdate, creation_date, update_date, deletion_date 
	FROM Clients 
	WHERE user_name = ? OR email = ?`

	//Execute query
	row := DB.QueryRow(query, email)

	//Scan results into client struct
	err := row.Scan(
		&client.UserID,
		&client.LastName,
		&client.FirstName,
		&client.UserName,
		&client.Email,
		&client.Password,
		&client.Avatar,
		&client.BirthDate,
		&client.BirthDate,
		&client.CreationDate,
		&client.UpdateDate,
		&client.DeletionDate,
	)

	//Check if client was found
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &client, nil
}

// Create client for testing
func InsertSampleClient() {
	// Check if the sample client already exists
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM Clients WHERE email = ?)", "sample@example.com").Scan(&exists)
	if err != nil {
		log.Printf("Error checking if sample client exists: %v", err)
		return
	}

	// Only insert if it does not exist
	if !exists {
		_, err = DB.Exec(`
            INSERT INTO Clients (last_name, first_name, user_name, email, password)
            VALUES ('Doe', 'John', 'johndoe', 'sample@example.com', 'securepassword')
        `)
		if err != nil {
			log.Printf("Failed to insert sample client: %v", err)
		} else {
			log.Println("Sample client inserted successfully.")
		}
	} else {
		log.Println("Sample client already exists, skipping insertion.")
	}
}
