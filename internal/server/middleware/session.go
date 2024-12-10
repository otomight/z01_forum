package middleware

import (
	"context"
	"forum/internal/config"
	"forum/internal/database"
	"log"
	"net/http"
	"time"
)

// Log each request's method & path
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Check session + set user role
func SessionMiddleWare(next http.Handler) http.Handler {
	var ctx context.Context

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Retrieve session from cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		//Get session data from database
		session, err := database.GetSessionByID(cookie.Value)
		if err != nil || time.Now().After(session.Expiration) {
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(r.Context(), config.SessionKey, session)
		//Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
