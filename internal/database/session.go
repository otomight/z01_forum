package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/config"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// Session CRUD operations
func CreateSession(session *UserSession) error {
	var	s	config.SessionsTableKeys

	s = config.TableKeys.Sessions
	if session.ID == "" || session.UserID == 0 {
		return fmt.Errorf("invalid session data: session_id & user_id can't be empty")
	}
	log.Printf("Attempting to create session: sessionID=%s, userID=%d, expiration=%v, userRole=%s",
		session.ID, session.UserID, session.Expiration, session.UserRole)
	_, err := insertInto(InsertIntoQuery{
		Table: s.Sessions,
		Keys: []string{s.ID, s.UserID, s.Expiration, s.UserRole, s.UserName},
		Values: [][]any{{
			session.ID, session.UserID, session.Expiration,
			session.UserRole, session.UserName,
		}},
	})
	if err != nil {
		log.Printf("Error creating session for user %d: %v", session.UserID, err)
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func GetSessionWithKey(key string, value any) (*UserSession, error) {
	var	query	string
	var	s		config.SessionsTableKeys
	var	row		*sql.Row
	var	session	UserSession
	var	err		error

	s = config.TableKeys.Sessions
	query = `
		SELECT `+s.ID+`, `+s.UserID+`, `+s.Expiration+`, `+s.CreationDate+`,
			`+s.UpdateDate+`, `+s.UserRole+`, `+s.UserName+`
		FROM `+s.Sessions+`
		WHERE `+key+` = ?
	`
	row = DB.QueryRow(query, value)
	err = row.Scan(
		&session.ID, &session.UserID, &session.Expiration,
		&session.CreationDate, &session.UpdateDate,
		&session.UserRole, &session.UserName,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("session not found: no matching session in database")
	}
	if err != nil {
		return nil, err // Return any other unexpected errors
	}
	if session.Expiration.Before(time.Now()) {
		return nil, errors.New("session expired")
	}
	return &session, nil
}

func GetSessionByID(sessionID string) (*UserSession, error) {
	return GetSessionWithKey(config.TableKeys.Sessions.ID, sessionID)
}

func GetSessionByUserID(userId int) (*UserSession, error) {
	return GetSessionWithKey(config.TableKeys.Sessions.UserID, userId)
}


// Create user session and return session ID
func CreateUserSession(userID int, userRole, userName string) (string, error) {
	sessionID := GenerateSessionID()
	expiration := time.Now().Add(24 * time.Hour)

	session := &UserSession{
		ID:			sessionID,
		UserID:		userID,
		UserRole:	userRole,
		UserName:	userName,
		Expiration:	expiration,
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
	var	s	config.SessionsTableKeys

	s = config.TableKeys.Sessions
	query := `
		DELETE FROM `+s.Sessions+`
		WHERE `+s.ID+` = ?
	`
	_, err := DB.Exec(query, sessionID)

	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
