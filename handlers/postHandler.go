package handlers

import (
	"Forum/database"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	//Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	//Check request method in POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Retreieve UserID from context
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	//Decode request body intot a Post struct
	var post database.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request playload", http.StatusBadRequest)
		return
	}

	//Post validation options
	if post.Title == "" || post.Content == "" {
		http.Error(w, "Title and Content are require", http.StatusBadRequest)
		return
	}

	//set creationDate + default values
	post.AuthorID = userID
	post.CreationDate = time.Now()
	post.UpdateDate = post.CreationDate
	post.IsDeleted = false

	//call function to save the post in database
	if err := database.CreatePost(&post); err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	//Respond with created post
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
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
	userRole, roleOk := r.Context().Value(UserRoleKey).(string)
	userID, idOk := r.Context().Value(UserIDKey).(int)

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

	//Redirect to appropriate page after deletion
	if userRole == "administrator" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else if userRole == "moderator" {
		http.Redirect(w, r, "/moderator", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
