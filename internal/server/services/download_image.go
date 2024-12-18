package services

import (
	"fmt"
	"forum/internal/server/models"
	"net/http"
	"os"
	"strconv"
)

func DownloadImage(w http.ResponseWriter, form *models.CreatePostForm, path string, id int) {
	// Create the directory, we will probably have to make this it's own function later on if we want to add images to comments
	dirpath := "./data/images/" + path + "/" + strconv.Itoa(id)
	os.MkdirAll(dirpath, 0755)

	// Save the file
	dst, err := os.Create(dirpath + "/" + form.Image.FileHeader.Filename)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.ReadFrom(form.Image.File)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", form.Image.FileHeader.Filename)
}
