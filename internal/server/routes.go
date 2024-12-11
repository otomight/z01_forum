package server

import (
	"forum/internal/server/handlers"
	"forum/internal/server/handlers/posthandlers"
	"net/http"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//Routes
	mux.Handle("/", sessionMiddleWare(http.HandlerFunc(handlers.HomePageHandler)))
	mux.Handle("/login", sessionMiddleWare(http.HandlerFunc(handlers.LoginHandler)))
	mux.Handle("/register", sessionMiddleWare(http.HandlerFunc(handlers.RegisterHandler)))
	mux.Handle("/logout", sessionMiddleWare(http.HandlerFunc(handlers.LogOutHandler)))

	// Rendering post creation form
	mux.Handle("/post/create", sessionMiddleWare(http.HandlerFunc(posthandlers.CreatePostHandler)))
	mux.Handle("/post/view/", sessionMiddleWare(http.HandlerFunc(posthandlers.ViewPostHandler)))

	mux.Handle("/categories/", sessionMiddleWare(http.HandlerFunc(handlers.CategoryPostsPageHandler)))

	mux.Handle("/history/created", sessionMiddleWare(http.HandlerFunc(handlers.HistoryCreatedPageHandler)))
	mux.Handle("/history/liked", sessionMiddleWare(http.HandlerFunc(handlers.HistoryLikedPageHandler)))

	//Protect route with session middelware
	mux.Handle("/post/delete", sessionMiddleWare(http.HandlerFunc(posthandlers.DeletePostHandler)))

	//Add comment
	mux.Handle("/post/comment/", sessionMiddleWare(http.HandlerFunc(handlers.AddCommentHandler)))

	mux.Handle("/post/like", sessionMiddleWare(http.HandlerFunc(handlers.PostLikeHandler)))
	mux.Handle("/post/dislike", sessionMiddleWare(http.HandlerFunc(handlers.PostDisLikeHandler)))
	mux.Handle("/comment/like", sessionMiddleWare(http.HandlerFunc(handlers.CommentLikeHandler)))
	mux.Handle("/comment/dislike", sessionMiddleWare(http.HandlerFunc(handlers.CommentDislikeHandler)))

	// Google log
	mux.Handle("/auth/google/login", sessionMiddleWare(http.HandlerFunc(handlers.GoogleLoginHandler)))
	mux.Handle("/Auth/google/callback", sessionMiddleWare(http.HandlerFunc(handlers.GoogleCallBackHandler)))

	// Github log
	mux.Handle("/auth/github/login", sessionMiddleWare(http.HandlerFunc(handlers.GithubLoginHandler)))
	mux.Handle("/auth/github/callback", sessionMiddleWare(http.HandlerFunc(handlers.GithubCallBackHandler)))

	// Facebook log
	mux.Handle("/auth/facebook/login", sessionMiddleWare(http.HandlerFunc(handlers.FacebookLoginHandler)))
	mux.Handle("/auth/facebook/callback", sessionMiddleWare(http.HandlerFunc(handlers.FacebookCallBackHandler)))

	//Wrap mux with logging middleware
	wrappedMux := loggingMiddleware(mux)

	return wrappedMux
}
