package services

import (
	db "forum/internal/database"
	"forum/internal/server/models"
	"time"
)

func CreatePost(userId int, form models.CreatePostForm) (int64, error) {
	post := &db.Post{
		AuthorID:     userId,
		Title:        form.Title,
		Content:      form.Content,
		Category:     form.Category,
		Tags:         form.Tags,
		CreationDate: time.Now(),
		UpdateDate:   time.Now(),
		IsDeleted:    false,
	}
	id, err := db.NewPost(post)
	if err != nil {
		return 0, err
	}
	return id, nil
}
