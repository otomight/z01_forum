package database

import (
	"time"
)

type Client struct {
	UserID       int       `json:"user_id"`
	LastName     string    `json:"last_name"`
	FirstName    string    `json:"first_name"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Avatar       string    `json:"avatar"`
	BirthDate    time.Time `json:"birth_date"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"udate_date"`
	DeletionDate time.Time `json:"deleteion_date"`
}

type Session struct {
	SessionID    string    `json:"session_id"`
	UserID       int       `json:"user_id"`
	Expiration   time.Time `json:"expiration"`
	CreationDate time.Time `json:"creation_date"`
	UpdateDate   time.Time `json:"update_date"`
	DeletionDate time.Time `json:"deletion_date"`
	IsDeleted    bool      `json:"is_deleted"`
}

type Post struct {
	PostID       int        `json:"post_id"`
	AuthorID     int        `json:"author_id"`
	Title        string     `json:"title"`
	Category     string     `json:"category"`
	Content      string     `json:"content"`
	CreationDate time.Time  `json:"creation_date"`
	UpdateDate   time.Time  `json:"update_date"`
	DeletionDate *time.Time `json:"deletion_date"`
	IsDeleted    bool       `json:"is_deleted"`
}
