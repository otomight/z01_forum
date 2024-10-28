package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Client CRUD operations
func CreateClient(client *Client) error {
	_, err := DB.Exec(`
		INSERT INTO Clients (last_name, first_name, user_name, email, password, avatar, bith_date, user_role)
		VALUES 	(?, ?, ?, ?, ?, ?, ?)
	`, client.LastName, client.FirstName, client.UserName, client.Email, client.Password, client.Avatar, client.BirthDate, client.UserRole)
	return err
}

func GetClientByID(userID int) (*Client, error) {
	row := DB.QueryRow(`
		SELECT user_id, last_name, first_name, user_name, email, avatar, birth_date, user_role, creation_date, update_date, deletion_date
		FROM Clients WERE user_id = ?
	`, userID)
	var client Client
	err := row.Scan(&client.UserID, &client.LastName, &client.FirstName, &client.UserName, &client.Email,
		&client.BirthDate, &client.UserRole, &client.CreationDate, &client.UpdateDate, &client.DeletionDate)
	return &client, err
}

// Save new user to database
func SaveUser(userName, email, password, firstName, lastName, userRole string) error {
	query := `INSERT INTO Clients (user_name, email, password, first_name, last_name, user_role)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query, userName, email, password, firstName, lastName, userRole)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}

// Retrieve client from database by their Email
func GetClientByUsernameOrEmail(email string) (*Client, error) {
	var client Client

	//Query to find user by Email
	query := `SELECT user_id, last_name, first_name, email, password, avatar, birthdate, user_role, creation_date, update_date, deletion_date 
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
		&client.UserRole,
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

// Setting user role
func GetUserRole(userID int) (string, error) {
	var role string

	query := "SELECT user_role FROM Clients WHERE user_id = ?"
	err := DB.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

// Validate User credentials
func ValidateUserCredentials(username, password string) (Client, error) {
	var user Client

	//Get user by username/email
	query := `SELECT user_id, first_name, last_name, user_name, email, password, user_role
	FROM Clients WHERE user_name = ? OR email = ? `

	row := DB.QueryRow(query, username, username)

	//Scan results into user struct
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password, &user.UserRole)
	if err == sql.ErrNoRows {
		return user, nil
	} else if err != nil {
		return user, err
	}

	//Compare hashed/provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, nil
	}
	return user, nil
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
            INSERT INTO Clients (last_name, first_name, user_name, email, password, user_role)
            VALUES ('Doe', 'John', 'johndoe', 'sample@example.com', 'securepassword', 'administrator')
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
