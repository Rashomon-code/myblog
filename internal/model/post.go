package model

import "time"

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Post struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
