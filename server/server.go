package server

import (
	"Forum/handlers"
	"Forum/middleware"
	"net/http"
)

func InitializeServer() http.Handler {
	mux := http.NewServeMux()

	//Routes
	mux.Handle("/", middleware.SessionMiddleWare(http.HandlerFunc(handlers.HomePageHandler)))
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/logout", handlers.LogOutHandler)

	// Rendering psot creation form
	mux.Handle("/posts/create", middleware.SessionMiddleWare(http.HandlerFunc(handlers.CreatePostFormHandler)))

	//Protect route with session middelware
	mux.Handle("/api/posts/create", middleware.SessionMiddleWare(http.HandlerFunc(handlers.CreatePostHandler)))
	mux.Handle("/posts/delete", middleware.SessionMiddleWare(http.HandlerFunc(handlers.DeletePostHandler)))
	mux.Handle("/posts/edit", middleware.SessionMiddleWare(http.HandlerFunc(handlers.EditPostHandler)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}
