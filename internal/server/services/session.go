package services

import (
	"forum/internal/database"
	"log"
	"net/http"
	"time"
)

func	GetSession(r *http.Request) (*database.UserSession, error) {
	var session			*database.UserSession
	var sessionCookie	*http.Cookie
	var err				error

	sessionCookie, err = r.Cookie("session_id")
	if err != nil || sessionCookie.Value == "" {
		log.Println("No session cookie found: user is not logged in")
		return &database.UserSession{
			IsLoggedIn: false,
		}, err
	}
	// Session cookie exists: retrieve session from DB
	session, err = database.GetSessionByID(sessionCookie.Value)
	if err != nil || time.Now().After(session.Expiration) {
		log.Println("Session expired or invalid")
		return &database.UserSession{
			IsLoggedIn: false,
		}, err
	}
	// Session is valid, user is logged in
	session.IsLoggedIn = true
	return session, nil
}
