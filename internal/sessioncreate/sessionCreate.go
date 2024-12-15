package sessioncreate

import (
	"forum/internal/database"
	"log"
	"net/http"
	"time"
)

func SessionCreate(
	w http.ResponseWriter, userID int, userRole string, userName string,
) error {
	sessionID, err := database.CreateUserSession(userID, userRole, userName)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Failed to create sesion", http.StatusInternalServerError)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
