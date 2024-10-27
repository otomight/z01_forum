package database

import "time"

type Client struct {
	UserID       int       `json:"user_id"`
	LastName     string    `json:"last_name"`
	FirstName    string    `json:"first_name"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Avatar       string    `json:"avatar"`
	BirthDate    time.Time `json:"birth_date"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"udate_date"`
	DeletionDate time.Time `json:"deleteion_date"`
}

type Session struct {
	SessionID    string    `json:"session_id"`
	UserID       int       `json:"user_id"`
	Expiration   time.Time `json:"expiration"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"update_date"`
	DeletionDate time.Time `json:"deletion_date"`
	IsDeleted    bool      `json:"is_deleted"`
}

type Post struct {
	PostID       int       `json:"post_id"`
	AuthorID     int       `json:"author_id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	Content      string    `json:"content"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"update_date"`
	DeletionDate time.Time `json:"deletion_date"`
	IsDeleted    bool      `json:"is_deleted"`
}

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

// Session CRUD operations
func CreateSession(session *Session) error {
	_, err := DB.Exec(`
		INSERT INTO Sessions (session_id, user_id, expiration)
		VALUES (?, ?, ?)
	`, session.SessionID, session.UserID, session.Expiration)
	return err
}

func GetSession(sessionID string) (*Session, error) {
	row := DB.QueryRow(`
		SELECT session_id, user_id, expiration, creation_date, update_date, deletion_date, is_deleted
		FROM Sessions WHERE session_id = ?
	`, sessionID)
	var session Session
	err := row.Scan(&session.SessionID, &session.UserID, &session.Expiration, &session.CreationDate,
		&session.UpdateDate, &session.DeletionDate, &session.IsDeleted)
	return &session, err
}
