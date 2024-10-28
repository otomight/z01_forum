package handlers

import "net/http"

func ModeratorDashboardHAndler(w http.ResponseWriter, r *http.Request) {
	userRole, ok := r.Context().Value(UserRoleKey).(string)
	if !ok || userRole != "moderator" {
		http.Error(w, "Acces denied", http.StatusForbidden)
		return
	}
}
