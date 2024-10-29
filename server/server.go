package server

import (
	"Forum/handlers"
	"Forum/middleware"
	"net/http"
)

func InitializeServer() http.Handler {
	mux := http.NewServeMux()

	//Routes
	mux.HandleFunc("/", handlers.RenderBaseHomePage)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)

	//Role redirections
	mux.Handle("/home", middleware.SessionMiddleWare(http.HandlerFunc(handlers.RenderBaseHomePage)))

	//Protect route with session middelware
	mux.Handle("posts/create", middleware.SessionMiddleWare(http.HandlerFunc(handlers.CreatePostHandler)))
	mux.Handle("/posts/delete", middleware.SessionMiddleWare(http.HandlerFunc(handlers.DeletePostHandler)))
	mux.Handle("/posts/edit", middleware.SessionMiddleWare(http.HandlerFunc(handlers.EditPostHandler)))

	//User dashboard
	mux.Handle("/user", middleware.SessionMiddleWare(http.HandlerFunc(handlers.UserHomePageHandler)))

	//Admin + moderator dashboards
	mux.Handle("/admin", middleware.SessionMiddleWare(http.HandlerFunc(handlers.AdministratorHomePageHandler)))
	mux.Handle("/moderator", middleware.SessionMiddleWare(http.HandlerFunc(handlers.ModeratorHomePageHandler)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}
