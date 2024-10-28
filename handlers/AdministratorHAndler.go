package handlers

import "net/http"

func AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(UserRoleKey).(string)
	if !ok || userRole != "administrator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}
}
