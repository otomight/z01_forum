package servsetup

import (
	"forum/internal/config"
	"forum/internal/server/handlers"
	"forum/internal/server/handlers/posthandlers"
	"log"
	"net/http"
)

func setupRoutes() *http.Handler {
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
	mux.Handle("/auth/callback", sessionMiddleWare(http.HandlerFunc(handlers.GoogleCallBackHandler)))

	// Discord log
	mux.Handle("/auth/discord/login", sessionMiddleWare(http.HandlerFunc(handlers.DiscordLoginHandler)))
	mux.Handle("/auth/discord/callback", sessionMiddleWare(http.HandlerFunc(handlers.DiscordCallBackHandler)))

	// Facebook log
	mux.Handle("/auth/facebook/login", sessionMiddleWare(http.HandlerFunc(handlers.FacebookLoginHandler)))
	mux.Handle("/auth/facebook/callback", sessionMiddleWare(http.HandlerFunc(handlers.FacebookCallBackHandler)))

	//Wrap mux with logging middleware
	wrappedMux := loggingMiddleware(mux)

	return &wrappedMux
}

func RedirectHTTP() {
	var	err	error

	err = http.ListenAndServe(
		":80", http.RedirectHandler(
			"https://localhost", http.StatusMovedPermanently,
		),
	)
	if err != nil {
		log.Fatalf("Could not start HTTP server for redirection: %v", err)
	}
}

func LaunchServer() {
	var	mux		*http.Handler
	var	server	*http.Server
	var	err		error

	mux = setupRoutes()
	server = &http.Server{
		Addr:		":https",
		Handler:	*mux,
	}
	log.Println("Starting server at https://localhost")
	err = server.ListenAndServeTLS(
		config.ServerCertifFilePath, config.ServerKeyFilePath,
	)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	go RedirectHTTP()
}
