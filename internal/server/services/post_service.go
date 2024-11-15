package services

import (
	db "forum/internal/database"
	"time"
)

func CreatePost(userId int, title string, content string,
							category string, tags string) (int64, error) {
	post := &db.Post{
		AuthorID:     userId,
		Title:        title,
		Content:      content,
		Category:     category,
		Tags:         tags,
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
