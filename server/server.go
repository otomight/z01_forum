package server

import (
	"Forum/handlers"
	"Forum/middleware"
	"net/http"
)

func InitializeServer() http.Handler {
	mux := http.NewServeMux()

	//Routes
	mux.HandleFunc("/posts", handlers.GetPost)

	//Protect route to create new Post with authentication
	mux.Handle("psots/create", middleware.LoggingMiddleware(http.HandlerFunc(handlers.CreatePost)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}
