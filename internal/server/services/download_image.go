package services

import (
	"forum/internal/config"
	"forum/internal/server/models"
	"forum/internal/utils"
	"log"
	"net/http"
	"os"
)

// Download the image and returns the image path
func DownloadImage(
	w http.ResponseWriter, form *models.CreatePostForm, dirPath string,
) string {
	var	imagePath		string
	var	imageServerPath	string
	var	dst				*os.File
	var	err				error

	imagePath = dirPath + form.Image.FileHeader.Filename
	imageServerPath, err = utils.CutStrAtPattern(imagePath, config.ImagesRoute)
	if err != nil {
		log.Println("Error at fetching the server version of the image path.")
		return ""
	}
	// Save the file
	dst, err = os.Create(imagePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return ""
	}
	defer dst.Close()

	_, err = dst.ReadFrom(form.Image.File)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return ""
	}
	log.Printf(
		"File uploaded successfully: %s\n", form.Image.FileHeader.Filename,
	)
	return imageServerPath
}
