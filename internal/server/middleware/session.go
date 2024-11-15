package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"forum/internal/database"
	"net/http"
	"time"
)

var sessionName = "user_session"

// Initialize gob encoder for USerSession struct
func init() {
	gob.Register(database.UserSession{})
}

// Save user session data in cookie
func SetSession(w http.ResponseWriter, r *http.Request, session database.UserSession) error {
	cookie := &http.Cookie{
		Name:     sessionName,
		Value:    encodeSession(session),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}

// Encode USerSession into a string
func encodeSession(session database.UserSession) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(session)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// Retrieve user session data from cookie
func GetSession(r *http.Request) (database.UserSession, error) {
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		return database.UserSession{}, err
	}
	return decodeSession(cookie.Value)
}

// Remove user session by expiring cookie
func DeleteSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1, //spontaneous
	}
	http.SetCookie(w, cookie)
}

// Decode string back into UserSession
func decodeSession(value string) (database.UserSession, error) {
	var session database.UserSession
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return session, err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&session)
	return session, err
}
