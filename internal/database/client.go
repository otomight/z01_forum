package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Client CRUD operations
func CreateClient(client *Client) error {
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		INSERT INTO `+k.Clients.Table+` (last_name, first_name, user_name,
						email, password, avatar, bith_date, user_role)
		VALUES 	(?, ?, ?, ?, ?, ?, ?)
	`
	_, err := DB.Exec(query, client.LastName, client.FirstName,
			client.UserName, client.Email, client.Password,
			client.Avatar, client.BirthDate, client.UserRole)
	return err
}

func GetClientByID(userID int) (*Client, error) {
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		SELECT user_id, last_name, first_name, user_name, email, avatar,
				birth_date, user_role, creation_date, update_date, deletion_date
		FROM `+k.Clients.Table+` WERE user_id = ?
	`
	row := DB.QueryRow(query, userID)
	var client Client
	err := row.Scan(&client.UserID, &client.LastName, &client.FirstName, &client.UserName, &client.Email,
		&client.BirthDate, &client.UserRole, &client.CreationDate, &client.UpdateDate, &client.DeletionDate)
	return &client, err
}

// Save new user to database
func SaveUser(
	userName, email, password string,
	firstName, lastName, userRole string,
) (int, error) {
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		INSERT INTO `+k.Clients.Table+` (user_name, email, password,
										first_name, last_name, user_role)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := DB.Exec(query, userName, email,
							password, firstName, lastName, userRole)
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
	var	k	config.StructTablesKeys = config.TableKeys

	//Query to find user by Email
	query := `
		SELECT user_id, last_name, first_name, email, password, avatar,
			birthdate, user_role, creation_date, update_date, deletion_date
		FROM `+k.Clients.Table+`
		WHERE user_name = ? OR email = ?
	`

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
	var	k	config.StructTablesKeys = config.TableKeys

	query := `
		SELECT user_role FROM `+k.Clients.Table+` WHERE user_id = ?
	`
	err := DB.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

// Validate User credentials
func ValidateUserCredentials(username, password string) (Client, error) {
	var user Client
	var	k	config.StructTablesKeys = config.TableKeys

	//Get user by username/email
	query := `
		SELECT user_id, first_name,
			last_name, user_name, email, password, user_role
		FROM `+k.Clients.Table+` WHERE user_name = ? OR email = ?
	`

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
	var user	Client
	var	k		config.StructTablesKeys = config.TableKeys
	// Check if the user already exists

	query := `
		SELECT user_id, last_name, first_name, user_name, email, avatar, user_role, oauth_provider, oauth_id
		FROM `+k.Clients.Table+`
		WHERE oauth_provider = ? AND oauth_id = ?
	`
	err := DB.QueryRow(query, oauthProvider, oauthID).Scan(
		&user.UserID, &user.LastName, &user.FirstName, &user.UserName, &user.Email,
		&user.Avatar, &user.UserRole, &user.OauthProvider, &user.OauthID,
	)

	// If no rows are found, create a new user
	if err == sql.ErrNoRows {
		// Insert new user
		query := `
			INSERT INTO `+k.Clients.Table+` (last_name, first_name, user_name,
				email, avatar, user_role, oauth_provider, oauth_id)
			VALUES (?, ?, ?, ?, ?, 'user', ?, ?)
		`
		result, insertErr := DB.Exec(query, "", name, name, email,
										avatar, oauthProvider, oauthID)
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
