package posthandlers

import (
	"fmt"
	"forum/internal/config"
	db "forum/internal/database"
	"forum/internal/server/models"
	"forum/internal/server/services"
	"forum/internal/server/templates"
	"forum/internal/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func createPost(userID int, form models.CreatePostForm) (int, error) {
	var err				error
	var categoriesIDs	[]int

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
	var	path	string
	var	err		error

	if err = utils.ParseForm(r, &form); err != nil {
		http.Error(w, "Unable to parse form:"+err.Error(),
			http.StatusBadRequest)
		return 0, err
	}
	if utils.IsStrEmpty(form.Title) ||
		utils.IsStrEmpty(form.Content) {
		http.Error(w, "Title and Content are required", http.StatusBadRequest)
		return 0, fmt.Errorf("Title and Content are required")
	}
	if len(form.Title) > config.PostTitleMaxLength {
		http.Error(w, "The title is too long", http.StatusBadRequest)
		return 0, fmt.Errorf("The title is too long")
	}
	if postID, err = createPost(session.UserID, form); err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to create post",
			http.StatusInternalServerError)
		return 0, err
	}
	if form.Image != nil {
		path = config.PostsImagesDirPath + strconv.Itoa(postID)
		os.MkdirAll(path, 0755)
		services.DownloadImage(w, &form, path)
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
			Session:    session,
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
