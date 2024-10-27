package middleware

import (
	"Forum/database"
	"Forum/handlers"
	"context"
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

// Check user's authentification
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Extract token from authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		//Validate session in database
		session, err := database.GetSessionByID(token)
		if err != nil {
			log.Printf("Error fetching session: %v", err)
			http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
			return
		}

		//Check session's expiration
		if session.IsDeleted || session.Expiration.Before(r.Context().Value("time").(time.Time)) {
			http.Error(w, "Unauthorized: session expired", http.StatusUnauthorized)
			return
		}

		//Attach userID to context
		ctx := context.WithValue(r.Context(), handlers.UserIDKey, session.UserID)
		r = r.WithContext(ctx)

		//Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
