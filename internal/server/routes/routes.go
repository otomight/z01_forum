package routes

import (
	"forum/internal/server/handlers"
	"forum/internal/server/middleware"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	//Routes
	mux.Handle("/", middleware.SessionMiddleWare(http.HandlerFunc(handlers.HomePageHandler)))
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/logout", handlers.LogOutHandler)

	// Rendering post creation form
	mux.Handle("/post/create", middleware.SessionMiddleWare(http.HandlerFunc(handlers.CreatePostHandler)))
	mux.Handle("/post/view/", middleware.SessionMiddleWare(http.HandlerFunc(handlers.DisplayPostHandler)))

	//Protect route with session middelware
	mux.Handle("/post/delete", middleware.SessionMiddleWare(http.HandlerFunc(handlers.DeletePostHandler)))
	mux.Handle("/post/edit", middleware.SessionMiddleWare(http.HandlerFunc(handlers.EditPostHandler)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}
