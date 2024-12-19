package servsetup

import (
	"crypto/tls"
	"forum/internal/config"
	"forum/internal/server/handlers"
	"forum/internal/server/handlers/posthandlers"
	"log"
	"net/http"
	"time"
)

func setupRoutes() *http.ServeMux {
	var	mux		*http.ServeMux

	mux = http.NewServeMux()

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
	mux.Handle("/auth/google/callback", sessionMiddleWare(http.HandlerFunc(handlers.GoogleCallBackHandler)))

	// Discord log
	mux.Handle("/auth/discord/login", sessionMiddleWare(http.HandlerFunc(handlers.DiscordLoginHandler)))
	mux.Handle("/auth/discord/callback", sessionMiddleWare(http.HandlerFunc(handlers.DiscordCallBackHandler)))

	// Facebook log
	mux.Handle("/auth/facebook/login", sessionMiddleWare(http.HandlerFunc(handlers.FacebookLoginHandler)))
	mux.Handle("/auth/facebook/callback", sessionMiddleWare(http.HandlerFunc(handlers.FacebookCallBackHandler)))

	//Wrap mux with logging middleware

	return mux
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

func getTLSConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates:	[]tls.Certificate{cert},
		MinVersion:		tls.VersionTLS13,
		MaxVersion:		tls.VersionTLS13,
		CipherSuites:	[]uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		PreferServerCipherSuites:	true,
		CurvePreferences:	[]tls.CurveID{
			tls.X25519,
		},
	}
}

func getServer(handler http.Handler) *http.Server {
	var	server	*http.Server
	var	tlsConf	*tls.Config
	var	cert	tls.Certificate
	var	err		error

	cert, err = tls.LoadX509KeyPair(
		config.ServerCertifFilePath, config.ServerKeyFilePath,
	)
	if err != nil {
		log.Fatalf("Error loading x509 key pair: %v", err)
	}
	tlsConf = getTLSConfig(cert)
	server = &http.Server{
		Addr:			":443",
		Handler:		handler,
		TLSConfig:		tlsConf,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
		IdleTimeout:	10 * time.Second,
	}
	return server
}

func LaunchServer() {
	var	mux			*http.ServeMux
	var	handler		http.Handler
	var	server		*http.Server
	var	rateLimiter	*RateLimiter
	var	err			error

	rateLimiter = newRateLimiter(500, time.Minute)
	mux = setupRoutes()
	handler = loggingMiddleware(mux)
	handler = rateLimiterMiddleware(rateLimiter, handler)
	server = getServer(handler)
	log.Println("Starting server at https://localhost")
	err = server.ListenAndServeTLS(
		config.ServerCertifFilePath, config.ServerKeyFilePath,
	)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
	go RedirectHTTP()
}
