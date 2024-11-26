package posthandlers

import (
	"fmt"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"net/http"
	"time"
)

func createPost(userId int, form models.CreatePostForm) (int64, error) {
	post := &db.Post{
		AuthorID:		userId,
		Title:			form.Title,
		Content:		form.Content,
		Category:		form.Category,
		Tags:			form.Tags,
		CreationDate:	time.Now(),
		UpdateDate:		time.Now(),
		IsDeleted:		false,
	}
	id, err := db.NewPost(post)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func createPostFromForm(
	w http.ResponseWriter, r *http.Request, session *db.UserSession,
) (int64, error) {
	var	form	models.CreatePostForm
	var	postId	int64
	var	err		error

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
	if postId, err = createPost(session.UserID, form); err != nil {
		http.Error(w, "Failed to create post",
							http.StatusInternalServerError)
		return 0, err
	}
	return postId, nil
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var	postId			int64
	var	redirectLink	string
	var	session			*db.UserSession
	var	data			models.CreatePostPageData
	var	ok				bool
	var	err				error

	session, ok = r.Context().Value(config.SessionKey).(*db.UserSession)
	if !ok {
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
