package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	var postIdStr	string
	var postId		int
	var post		*database.Post
	var data		models.ViewPostPageData
	var session		*database.UserSession
	var comments	[]database.Comment
	var err			error

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	postIdStr = strings.TrimPrefix(r.URL.Path, "/post/view/")
	if postIdStr == "" || strings.Contains(postIdStr, "/") {
		http.NotFound(w, r)
		return
	}
	postId, err = strconv.Atoi(postIdStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	post, err = database.GetPostByID(postId)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	//Fetch comments for the post
	comments, err = database.GetCommentsByPostID(postId)
	if err != nil {
		log.Printf("Failed to fetch comments for post %d: %v", postId, err)
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}
	session, _ = services.GetSession(r)
	data = models.ViewPostPageData{
		Post: post,
		Session: session,
		Comments: comments,
	}
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}

func createPostFromForm(w http.ResponseWriter,
							r *http.Request) (int64, error) {
	var ok		bool
	var userId	int
	var form	models.CreatePostForm
	var postId	int64
	var err		error

	userId, ok = r.Context().Value(config.UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context",
						http.StatusInternalServerError)
		return 0, err
	}
	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return 0, err
	}
	if form.Title == "" || form.Content == "" {
		http.Error(w, "Title and Content are required",
			http.StatusBadRequest)
		return 0, err
	}
	if postId, err = services.CreatePost(userId, form); err != nil {
		http.Error(w, "Failed to create post",
			http.StatusInternalServerError)
		return 0, err
	}
	return postId, nil
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var postId			int64
	var redirectLink	string
	var session			*database.UserSession
	var data			models.CreatePostPageData
	var err				error

	if r.Method == http.MethodGet {
		// render the post creation page
		session, err = services.GetSession(r)
		if err != nil {
			http.Error(w, "User not logged", http.StatusUnauthorized)
			return
		}
		data = models.CreatePostPageData{
			Session: session,
		}
		templates.RenderTemplate(w, config.CreatePostTmpl, data)
	} else if r.Method == http.MethodPost {
		// handle the form send on post creation
		if postId, err = createPostFromForm(w, r); err != nil {
			return
		}
		redirectLink = fmt.Sprintf("/post/view/%d", postId)
		http.Redirect(w, r, redirectLink, http.StatusSeeOther)
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	posts, err := database.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userRole, roleOk := r.Context().Value(config.UserRoleKey).(string)
	userID, idOk := r.Context().Value(config.UserIDKey).(int)

	//Check authentication
	if !roleOk || !idOk {
		http.Error(w, "Unauthorizes", http.StatusUnauthorized)
		return
	}

	//Retrieve postID from URL parameters
	postIDstr := r.URL.Query().Get("postID")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil || postID == 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	//Fetch post to identify its author
	post, err := database.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	//Authorization check
	if userRole == "administrator" || userRole == "moderator" || post.AuthorID == userID {
		//Delete post if authorized
		err = database.DeletePost(postID)
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unauthorized to delete this post", http.StatusInternalServerError)
		return
	}
}

// NOT FINISHED
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, idOk := r.Context().Value(config.UserIDKey).(int)
	userRole, _ := r.Context().Value(config.UserRoleKey).(string)

	//Check authentication
	if !idOk {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Retrieve postID from URL parameters
	postIDstr := r.URL.Query().Get("postID")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil || postID == 0 {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	//Fetch post to identify its author
	post, err := database.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	//Only Authors can edit their posts
	if post.AuthorID != userID {
		http.Error(w, "Unauthorized to edit this post", http.StatusForbidden)
		return
	}
	data := struct {
		Title      string
		Post       *database.Post
		IsLoggedIn bool
		UserRole   string
	}{
		Title:      "Edit Post",
		Post:       post,
		IsLoggedIn: true,
		UserRole:   userRole,
	}
	templates.RenderTemplate(w, config.EditPostTmpl, data)
}
