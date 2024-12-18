package posthandlers

import (
	"fmt"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"time"
)

func createPost(userID int, form models.CreatePostForm) (int, error) {
	var	err				error
	var	categoriesIDs	[]int

	post := &db.Post{
		AuthorID:		userID,
		Title:			form.Title,
		Content:		form.Content,
		CreationDate:	time.Now(),
		UpdateDate:		time.Now(),
	}
	categoriesIDs, err = utils.StrSliceToIntSlice(form.Categories)
	if err != nil {
		log.Printf("Error during convertion of categories IDs: %v", err)
	}
	id, err := db.NewPost(post, categoriesIDs)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func createPostFromForm(
	w http.ResponseWriter, r *http.Request, session *db.UserSession,
) (int, error) {
	var	form	models.CreatePostForm
	var	postID	int
	var	err		error

	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:" + err.Error(),
							http.StatusBadRequest)
		return 0, err
	}
	if !utils.IsOnlyPrintable(form.Title) ||
	!utils.IsOnlyPrintable(form.Content) {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return 0, err
	}
	if postID, err = createPost(session.UserID, form); err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to create post",
							http.StatusInternalServerError)
		return 0, err
	}
	return postID, nil
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var	postID			int
	var	redirectLink	string
	var	session			*db.UserSession
	var	categories		[]*db.Category
	var	data			models.CreatePostPageData
	var	err				error

	session, _ = r.Context().Value(config.SessionKey).(*db.UserSession)
	if session == nil {
		http.Error(w, "User not logged", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
	// render the post creation page
		if categories, err = db.GetGlobalCategories(); err != nil {
			http.Error(
				w, "Error at fetching categories",
				http.StatusInternalServerError,
			)
		}
		data = models.CreatePostPageData{
			Session: session,
			Categories: categories,
		}
		templates.RenderTemplate(w, config.CreatePostTmpl, data)
	} else if r.Method == http.MethodPost {
	// handle the form send on post creation
		if postID, err = createPostFromForm(w, r, session); err != nil {
			return
		}
		redirectLink = fmt.Sprintf("/post/view/%d", postID)
		http.Redirect(w, r, redirectLink, http.StatusSeeOther)
	}
}
