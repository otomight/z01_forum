package handlers

import (
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"log"
	"net/http"
)

func	IncludeUserLikesConfigs(
	session *db.UserSession, posts []db.Post,
) []*models.PostWithUserConfig {
	var	userPost	*models.PostWithUserConfig
	var	userPosts	[]*models.PostWithUserConfig
	var	i			int

	for i = 0; i < len(posts); i++ {
		userPost = &models.PostWithUserConfig{}
		userPost.Post = &posts[i]
		userPost.IsLikedByUser, userPost.IsDislikedByUser =
						services.GetUserLikesConfigsOfPost(session, &posts[i])
		userPosts = append(userPosts, userPost)
	}
	return userPosts
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	var	session		*db.UserSession
	var	posts		[]db.Post
	var	userPosts	[]*models.PostWithUserConfig
	var	err			error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	posts, err = db.GetAllPosts()
	// Prepare posts
	if err != nil {
		log.Printf("Failed to retrieve posts: %v", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	userPosts = IncludeUserLikesConfigs(session, posts)
	data := models.HomePageData{
		Posts:		userPosts,
		Session:	session,
	}
	templates.RenderTemplate(w, config.HomeTmpl, data)
}
