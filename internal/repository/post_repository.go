package repository

import (
	"database/sql"
	"errors"
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
		VALUES ($1, $2, $3)	
	`
	_, err := r.db.Exec(insertSQL, userID, title, content)
	if err != nil {
		return fmt.Errorf("投稿失敗: %w", err)
	}

	return nil
}

func (r *PostRepository) GetTitleByUserID(userID int64) ([]model.ArticleSummary, error) {
	selectSQL := `SELECT id, title, created_at FROM posts WHERE user_id = $1`

	rows, err := r.db.Query(selectSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleSummaries(rows)
}

func (r *PostRepository) GetPostDetail(postID int64) (model.PostDetail, error) {
	selectSQL := `
		SELECT p.id, p.title, p.content, p.user_id, p.created_at, up.display_name
		FROM posts p
		LEFT JOIN user_profiles up ON p.user_id = up.user_id
		WHERE p.id = $1
	`

	var p model.PostDetail
	row := r.db.QueryRow(selectSQL, postID)
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CreatedAt, &p.DisplayName)
	if err != nil {
		return model.PostDetail{}, fmt.Errorf("文章が読み取れませんでした: %w", err)
	}

	if p.DisplayName == nil || *p.DisplayName == "" {
		displayName := fmt.Sprintf("ユーザー%d", p.UserID)
		p.DisplayName = &displayName
	}

	return p, nil
}

func (r *PostRepository) DeletePost(postID int64) error {
	deleteSQL := `DELETE FROM posts WHERE id = $1`

	result, err := r.db.Exec(deleteSQL, postID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("何も削除していませんでした")
	}

	return nil
}

func (r *PostRepository) EditPost(postID int64, title string, content string) error {
	updateSQL := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`

	result, err := r.db.Exec(updateSQL, title, content, postID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("更新できませんでした")
	}

	return nil
}

func scanArticleSummaries(rows *sql.Rows) ([]model.ArticleSummary, error) {
	defer rows.Close()

	var posts []model.ArticleSummary
	for rows.Next() {
		var a model.ArticleSummary
		if err := rows.Scan(&a.ID, &a.Title, &a.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetAllPost() ([]model.ArticleSummary, error) {
	selectSQL := `SELECT id, title, created_at FROM posts ORDER BY created_at DESC`

	rows, err := r.db.Query(selectSQL)
	if err != nil {
		return nil, err
	}

	return scanArticleSummaries(rows)
}

func (r *PostRepository) SearchPost(keyword string) ([]model.ArticleSummary, error) {
	selectSQL := `SELECT id, title, created_at FROM posts WHERE title LIKE '%' || $1 || '%' ORDER BY created_at DESC`

	rows, err := r.db.Query(selectSQL, keyword)
	if err != nil {
		return nil, err
	}
	posts, err := scanArticleSummaries(rows)

	return posts, nil
}
