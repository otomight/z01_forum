package handlers

import (
	"Forum/database"
	"encoding/json"
	"net/http"
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
