package services

import (
	"fmt"
	"forum/internal/server/models"
	"net/http"
	"os"
)

func DownloadImage(w http.ResponseWriter, form *models.CreatePostForm, path string) {
	// Save the file
	dst, err := os.Create(path + "/" + form.Image.FileHeader.Filename)
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
