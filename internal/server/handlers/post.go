package handlers

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"forum/internal/utils"
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
	if postId, err = strconv.Atoi(postIdStr); err != nil {
		http.NotFound(w, r)
		return
	}
	if post, err = database.GetPostByID(postId); err != nil {
		http.NotFound(w, r)
		return
	}
	session, _ = services.GetSession(r)
	data = models.ViewPostPageData{
		Post: post,
		Session: session,
	}
	templates.RenderTemplate(w, config.ViewPostTmpl, data)
}

func createPostFromForm(w http.ResponseWriter,
			r *http.Request, session *database.UserSession) (int64, error) {
	var form	models.CreatePostForm
	var postId	int64
	var err		error

	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:" + err.Error(),
					http.StatusBadRequest)
		return 0, err
	}
	if form.Title == "" || form.Content == "" {
		http.Error(w, "Title and Content are required",
					http.StatusBadRequest)
		return 0, err
	}
	if postId, err = services.CreatePost(session.UserID, form); err != nil {
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

	session, err = services.GetSession(r)
	if err != nil {
		http.Error(w, "User not logged", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		// render the post creation page
		data = models.CreatePostPageData{
			Session: session,
		}
		templates.RenderTemplate(w, config.CreatePostTmpl, data)
	} else if r.Method == http.MethodPost {
		// handle the form send on post creation
		if postId, err = createPostFromForm(w, r, session); err != nil {
			return
		}
		redirectLink = fmt.Sprintf("/post/view/%d", postId)
		http.Redirect(w, r, redirectLink, http.StatusSeeOther)
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	var form			models.DeletePostForm
	var err				error
	var postId			int

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return
	}
	if postId, err = strconv.Atoi(form.PostId); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	if err = database.DeletePost(postId); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {

}
