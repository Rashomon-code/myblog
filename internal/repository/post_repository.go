package repository

import (
	"database/sql"
	"fmt"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(userID int64, title string, content string) error {
	insertSQL := `
		INSERT INTO posts (user_id, title, content)
		VALUES (?, ?, ?)	
	`
	_, err := r.db.Exec(insertSQL, userID, title, content)
	if err != nil {
		return fmt.Errorf("投稿失敗: %w", err)
	}

	fmt.Println("投稿しました。")
	return nil
}
