package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// Session CRUD operations
func CreateSession(session *Session) error {
	if session.SessionID == "" || session.UserID == 0 {
		return fmt.Errorf("invalid session data: session_id & user_id can't be empty")
	}
	_, err := DB.Exec(`
		INSERT INTO Sessions (session_id, user_id, expiration)
		VALUES (?, ?, ?)
	`, session.SessionID, session.UserID, session.Expiration)

	if err != nil {
		log.Printf("Error creatin session for user %d: %v", session.UserID, err)
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func GetSessionByID(sessionID string) (*Session, error) {
	row := DB.QueryRow(`
		SELECT session_id, user_id, expiration, creation_date, update_date, deletion_date, is_deleted
		FROM Sessions WHERE session_id = ?
	`, sessionID)
	var session Session
	err := row.Scan(&session.SessionID, &session.UserID, &session.Expiration, &session.CreationDate,
		&session.UpdateDate, &session.DeletionDate, &session.IsDeleted)
	return &session, err
}

// Create user session and return session ID
func CreateUserSession(userID int) (string, error) {
	sessionID := GenerateSessionID()
	expiration := time.Now().Add(24 * time.Hour)

	session := &Session{
		SessionID:  sessionID,
		UserID:     userID,
		Expiration: expiration,
	}

	//Store session to database
	if err := CreateSession(session); err != nil {
		return "", err
	}
	return sessionID, nil
}

// Generate unique session id
func GenerateSessionID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
