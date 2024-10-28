package server

import (
	"Forum/handlers"
	"Forum/middleware"
	"net/http"
)

func InitializeServer() http.Handler {
	mux := http.NewServeMux()

	//Routes
	mux.HandleFunc("/", handlers.RenderHomePage)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHAndler)
	mux.HandleFunc("/admin", handlers.AdminDashBoard)
	mux.HandleFunc("/moderator", handlers.ModeratorDashboardHAndler)
	mux.HandleFunc("/posts/delete", handlers.DeletePostHandler)

	//Protect route to create new Post with authentication
	mux.Handle("posts/create", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePost)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}
