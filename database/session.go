package database

import (
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
		INSERT INTO Sessions (session_id, user_id, expiration, user_role)
		VALUES (?, ?, ?, ?)
	`, session.SessionID, session.UserID, session.Expiration, session.UserRole)

	if err != nil {
		log.Printf("Error creating session for user %d: %v", session.UserID, err)
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func GetSessionByID(sessionID string) (*UserSession, error) {
	row := DB.QueryRow(`
		SELECT session_id, user_id, expiration, creation_date, update_date, deletion_date, is_deleted, user_role
		FROM Sessions WHERE session_id = ?
	`, sessionID)
	var session UserSession
	err := row.Scan(&session.SessionID, &session.UserID, &session.Expiration, &session.CreationDate,
		&session.UpdateDate, &session.DeletionDate, &session.IsDeleted, &session.UserRole)
	if err != nil {
		return nil, err
	}
	return &session, err
}

// Create user session and return session ID
func CreateUserSession(userID int, userRole string) (string, error) {
	sessionID := GenerateSessionID()
	expiration := time.Now().Add(24 * time.Hour)

	session := &UserSession{
		SessionID:  sessionID,
		UserID:     userID,
		UserRole:   userRole,
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
