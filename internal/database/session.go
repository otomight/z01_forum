package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// Session CRUD operations
func CreateSession(session *UserSession) error {
	if session.SessionID == "" || session.UserID == 0 {
		return fmt.Errorf("invalid session data: session_id & user_id can't be empty")
	}

	log.Printf("Attempting to create session: sessionID=%s, userID=%d, expiration=%v, userRole=%s",
		session.SessionID, session.UserID, session.Expiration, session.UserRole)

	_, err := DB.Exec(`
		INSERT INTO sessions (session_id, user_id, expiration, user_role, user_name)
		VALUES (?, ?, ?, ?, ?)
	`, session.SessionID, session.UserID, session.Expiration, session.UserRole, session.UserName)

	if err != nil {
		log.Printf("Error creating session for user %d: %v", session.UserID, err)
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func GetSessionByID(sessionID string) (*UserSession, error) {
	row := DB.QueryRow(`
		SELECT session_id, user_id, expiration, creation_date, update_date, deletion_date, is_deleted, user_role, user_name
		FROM sessions WHERE session_id = ?
	`, sessionID)
	var session UserSession
	err := row.Scan(&session.SessionID, &session.UserID, &session.Expiration, &session.CreationDate,
		&session.UpdateDate, &session.DeletionDate, &session.IsDeleted, &session.UserRole, &session.UserName)

	// Handle specific error if no rows are returned
	if err == sql.ErrNoRows {
		return nil, errors.New("session not found: no matching session in database")
	}
	if err != nil {
		return nil, err // Return any other unexpected errors
	}

	// Additional checks to see if the session is expired or marked as deleted
	if session.IsDeleted || session.Expiration.Before(time.Now()) {
		return nil, errors.New("session invalid: expired or marked as deleted")
	}

	return &session, nil
}

// Create user session and return session ID
func CreateUserSession(userID int, userRole, userName string) (string, error) {
	sessionID := GenerateSessionID()
	expiration := time.Now().Add(24 * time.Hour)

	session := &UserSession{
		SessionID:  sessionID,
		UserID:     userID,
		UserRole:   userRole,
		UserName:   userName,
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

func DeleteSession(sessionID string) error {
	_, err := DB.Exec(`
		DELETE FROM sessions
		WHERE session_id = ?
	`, sessionID)

	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
