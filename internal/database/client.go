package database

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Client CRUD operations
func CreateClient(client *Client) error {
	_, err := DB.Exec(`
		INSERT INTO clients (last_name, first_name, user_name, email, password, avatar, bith_date, user_role)
		VALUES 	(?, ?, ?, ?, ?, ?, ?)
	`, client.LastName, client.FirstName, client.UserName, client.Email, client.Password, client.Avatar, client.BirthDate, client.UserRole)
	return err
}

func GetClientByID(userID int) (*Client, error) {
	row := DB.QueryRow(`
		SELECT user_id, last_name, first_name, user_name, email, avatar, birth_date, user_role, creation_date, update_date, deletion_date
		FROM clients WERE user_id = ?
	`, userID)
	var client Client
	err := row.Scan(&client.UserID, &client.LastName, &client.FirstName, &client.UserName, &client.Email,
		&client.BirthDate, &client.UserRole, &client.CreationDate, &client.UpdateDate, &client.DeletionDate)
	return &client, err
}

// Save new user to database
func SaveUser(userName, email, password, firstName, lastName, userRole string) (int, error) {
	query := `
	INSERT INTO clients (user_name, email, password, first_name, last_name, user_role)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(query, userName, email, password, firstName, lastName, userRole)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return int(userID), nil
}

// Retrieve client from database by their Email
func GetClientByUsernameOrEmail(email string) (*Client, error) {
	var client Client

	//Query to find user by Email
	query := `SELECT user_id, last_name, first_name, email, password, avatar, birthdate, user_role, creation_date, update_date, deletion_date
	FROM clients
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

	query := "SELECT user_role FROM clients WHERE user_id = ?"
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
	FROM clients WHERE user_name = ? OR email = ? `

	row := DB.QueryRow(query, username, username)

	//Scan results into user struct
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password, &user.UserRole)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, err
	}

	//Compare hashed/provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, fmt.Errorf("invalid password")
	}
	return user, nil
}

// Social Login
func GetOrCreateUserByOAuth(oauthProvider, oauthID, email, name, avatar string) (*Client, error) {
	var user Client
	// Check if the user already exists
	err := DB.QueryRow(`
		SELECT user_id, last_name, first_name, user_name, email, avatar, user_role, oauth_provider, oauth_id
		FROM clients
		WHERE oauth_provider = ? AND oauth_id = ?
	`, oauthProvider, oauthID).Scan(
		&user.UserID, &user.LastName, &user.FirstName, &user.UserName, &user.Email,
		&user.Avatar, &user.UserRole, &user.OauthProvider, &user.OauthID,
	)

	// If no rows are found, create a new user
	if err == sql.ErrNoRows {
		// Insert new user
		result, insertErr := DB.Exec(`
			INSERT INTO clients (last_name, first_name, user_name, email, avatar, user_role, oauth_provider, oauth_id)
			VALUES (?, ?, ?, ?, ?, 'user', ?, ?)
		`, "", name, name, email, avatar, oauthProvider, oauthID)
		if insertErr != nil {
			log.Printf("Insert failed: %v", insertErr)
			return nil, fmt.Errorf("failed to create user: %w", insertErr)
		}

		// Get the ID of the newly inserted user
		userID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Failed to retrieve new userID:%v", err)
			return nil, fmt.Errorf("failed to retrieve new user ID: %w", err)
		}
		log.Printf("User inserted with ID: %d", userID)

		// Return the newly created user
		return &Client{
			UserID:        int(userID),
			FirstName:     name,
			UserName:      name,
			Email:         email,
			Avatar:        avatar,
			UserRole:      "user",
			OauthProvider: oauthProvider,
			OauthID:       oauthID,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}

// Create client for testing
func InsertSampleClient() {
	// Check if the sample client already exists
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clients WHERE email = ?)", "sample@example.com").Scan(&exists)
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
		_, err = DB.Exec(`
            INSERT INTO clients (last_name, first_name, user_name, email, password, user_role)
            VALUES ('Doe', 'John', 'johndoe', 'sample@example.com', ?, 'administrator')
        `, hashedPassword)
		if err != nil {
			log.Printf("Failed to insert sample client: %v", err)
		} else {
			log.Println("Sample client inserted successfully.")
		}
	} else {
		log.Println("Sample client already exists, skipping insertion.")
	}
}
