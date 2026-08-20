package model

import "time"

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ArticleSummary struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PostDetail struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName *string   `json:"display_name"` //LEFT JOIN で DisplayName が存在しない場合 nil になるため pointer を用いる
}
