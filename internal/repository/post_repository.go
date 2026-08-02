package repository

import (
	"database/sql"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/model"
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

func (r *PostRepository) GetTitleByUserID(userID int64) ([]model.ArticleSummary, error) {
	selectSQL := `SELECT id, title, created_at FROM posts WHERE user_id = ?`

	rows, err := r.db.Query(selectSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.ArticleSummary
	for rows.Next() {
		var a model.ArticleSummary
		err := rows.Scan(&a.ID, &a.Title, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostDetail(postID int64) ([]model.Post, error) {
	selectSQL := `SELECT tilte, content, user_id, created_at FROM posts WHERE id = ?`

	rows, err := r.db.Query(selectSQL, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.Title, &p.Content, &p.UserID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
