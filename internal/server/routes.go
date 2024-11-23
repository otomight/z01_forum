package server

import (
	"forum/internal/server/handlers"
	"forum/internal/server/handlers/posthandlers"
	"forum/internal/server/middleware"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//Routes
	mux.Handle("/", middleware.SessionMiddleWare(http.HandlerFunc(handlers.HomePageHandler)))
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/logout", handlers.LogOutHandler)

	// Rendering post creation form
	mux.Handle("/post/create", middleware.SessionMiddleWare(http.HandlerFunc(posthandlers.CreatePostHandler)))
	mux.Handle("/post/view/", middleware.SessionMiddleWare(http.HandlerFunc(posthandlers.ViewPostHandler)))

	//Protect route with session middelware
	mux.Handle("/post/delete", middleware.SessionMiddleWare(http.HandlerFunc(posthandlers.DeletePostHandler)))
	mux.Handle("/post/edit", middleware.SessionMiddleWare(http.HandlerFunc(posthandlers.EditPostHandler)))

	//Add comment
	mux.Handle("/post/comment/", middleware.SessionMiddleWare(http.HandlerFunc(handlers.AddCommentHandler)))

	mux.Handle("/post/like", middleware.SessionMiddleWare(http.HandlerFunc(handlers.LikePostHandler)))
	mux.Handle("/post/dislike", middleware.SessionMiddleWare(http.HandlerFunc(handlers.DisLikePostHandler)))

	//Wrap mux with logging middleware
	wrappedMux := middleware.LoggingMiddleware(mux)

	return wrappedMux
}