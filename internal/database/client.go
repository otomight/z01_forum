package database

import (
	"database/sql"
	"fmt"
	"forum/internal/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Client CRUD operations
// func CreateClient(client *Client) error {
// 	var	c	config.ClientsTableKeys

// 	c = config.TableKeys.Clients
// 	query := `
// 		INSERT INTO `+c.Clients+` (
// 			`+c.LastName+`, `+c.FirstName+`, `+c.UserName+`,
// 			`+c.Email+`, `+c.Password+`,
// 			`+c.Avatar+`, `+c.BirthDate+`, `+c.UserRole+`
// 		)
// 		VALUES 	(?, ?, ?, ?, ?, ?, ?)
// 	`
// 	_, err := DB.Exec(query, client.LastName, client.FirstName,
// 			client.UserName, client.Email, client.Password,
// 			client.Avatar, client.BirthDate, client.UserRole)
// 	return err
// }

// func GetClientByID(userID int) (*Client, error) {
// 	var	c	config.ClientsTableKeys

// 	c = config.TableKeys.Clients
// 	query := `
// 		SELECT `+c.ID+`, `+c.LastName+`, `+c.FirstName+`, `+c.UserName+`,
// 				`+c.Email+`, `+c.Avatar+`, `+c.BirthDate+`, `+c.UserRole+`,
// 				`+c.CreationDate+`, `+c.UpdateDate+`, `+c.DeletionDate+`
// 		FROM `+c.Clients+` WERE `+c.ID+` = ?
// 	`
// 	row := DB.QueryRow(query, userID)
// 	var client Client
// 	err := row.Scan(&client.ID, &client.LastName, &client.FirstName, &client.UserName, &client.Email,
// 		&client.BirthDate, &client.UserRole, &client.CreationDate, &client.UpdateDate, &client.DeletionDate)
// 	return &client, err
// }

// Save new user to database
func SaveUser(
	userName string, email string, password string, userRole string,
) (int, error) {
	var	c	config.ClientsTableKeys

	c = config.TableKeys.Clients
	result, err := insertInto(InsertIntoQuery{
		Table: c.Clients,
		Keys: []string{
			c.UserName, c.Email, c.Password, c.UserRole,
		},
		Values: [][]any{{
			userName, email, password, userRole,
		}},
	})
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
// func GetClientByUsernameOrEmail(email string) (*Client, error) {
// 	var	client	Client
// 	var	c		config.ClientsTableKeys

// 	//Query to find user by Email
// 	c = config.TableKeys.Clients
// 	query := `
// 		SELECT `+c.ID+`, `+c.LastName+`, `+c.FirstName+`, `+c.Email+`,
// 				`+c.Password+`, `+c.Avatar+`, `+c.BirthDate+`, `+c.UserRole+`,
// 				`+c.CreationDate+`, `+c.UpdateDate+`, `+c.DeletionDate+`
// 		FROM `+c.Clients+`
// 		WHERE `+c.UserName+` = ? OR `+c.Email+` = ?
// 	`

// 	//Execute query
// 	row := DB.QueryRow(query, email)

// 	//Scan results into client struct
// 	err := row.Scan(
// 		&client.ID,
// 		&client.LastName,
// 		&client.FirstName,
// 		&client.UserName,
// 		&client.Email,
// 		&client.Password,
// 		&client.Avatar,
// 		&client.BirthDate,
// 		&client.BirthDate,
// 		&client.UserRole,
// 		&client.CreationDate,
// 		&client.UpdateDate,
// 		&client.DeletionDate,
// 	)

// 	//Check if client was found
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return &client, nil
// }

// Setting user role
// func GetUserRole(userID int) (string, error) {
// 	var role string
// 	var	c	config.ClientsTableKeys

// 	c = config.TableKeys.Clients
// 	query := `
// 		SELECT `+c.UserRole+` FROM `+c.Clients+` WHERE `+c.UserName+` = ?
// 	`
// 	err := DB.QueryRow(query, userID).Scan(&role)
// 	if err != nil {
// 		return "", err
// 	}
// 	return role, nil
// }

// Validate User credentials
func ValidateUserCredentials(username, password string) (Client, error) {
	var user Client
	var	c	config.ClientsTableKeys

	c = config.TableKeys.Clients
	//Get user by username/email
	query := `
		SELECT `+c.ID+`, `+c.UserName+`, `+c.Email+`,
				`+c.Password+`, `+c.UserRole+`
		FROM `+c.Clients+` WHERE `+c.UserName+` = ? OR `+c.Email+` = ?;
	`

	row := DB.QueryRow(query, username, username)

	//Scan results into user struct
	err := row.Scan(
		&user.ID, &user.UserName, &user.Email, &user.Password, &user.UserRole,
	)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, err
	}

	//Compare hashed/provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("invalid password")
	}
	return user, nil
}

// Social Login
func GetOrCreateUserByOAuth(oauthProvider, oauthID, email, name, avatar string) (*Client, error) {
	var user	Client
	var	c		config.ClientsTableKeys

	c = config.TableKeys.Clients
	// Check if the user already exists
	query := `
		SELECT `+c.ID+`, `+c.LastName+`, `+c.FirstName+`,
				`+c.UserName+`, `+c.Email+`, `+c.Avatar+`,
				`+c.UserRole+`, `+c.OauthProvider+`, `+c.OauthID+`
		FROM `+c.Clients+`
		WHERE `+c.OauthProvider+` = ? AND `+c.OauthID+` = ?
	`
	err := DB.QueryRow(query, oauthProvider, oauthID).Scan(
		&user.ID, &user.LastName, &user.FirstName, &user.UserName, &user.Email,
		&user.Avatar, &user.UserRole, &user.OauthProvider, &user.OauthID,
	)

	// If no rows are found, create a new user
	if err == sql.ErrNoRows {
		// Insert new user
		query := `
			INSERT INTO `+c.Clients+` (
				`+c.LastName+`, `+c.FirstName+`, `+c.UserName+`,
				`+c.Email+`, `+c.Avatar+`, `+c.UserRole+`,
				`+c.OauthProvider+`, `+c.OauthID+`
			)
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
			ID:        int(userID),
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
