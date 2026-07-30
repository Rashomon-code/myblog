package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetPasswordHash(username string) (string, error) {
	selectSQL := `SELECT password_hash FROM users WHERE username = ?`

	var passwordHash string

	row := r.db.QueryRow(selectSQL, username)
	err := row.Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("入力に誤りがございます。")
		}

		return "", fmt.Errorf("予期せぬエラーが起きました: %w", err)
	}
	return passwordHash, nil
}

func (r *UserRepository) CreateUser(username, passwordHash string) error {
	insertSQL := `
		INSERT INTO users (username, password_hash)
		VALUES (?, ?)
	`

	_, err := r.db.Exec(insertSQL, username, passwordHash)
	if err != nil {
		return fmt.Errorf("登録にエラーが起きました: %w", err)
	}
	fmt.Println("登録完了しました。")
	return nil
}

func (r *UserRepository) FindByUsername(username string) (int64, error) {
	selectSQL := `SELECT id FROM users WHERE username = ?`

	var userID int64

	row := r.db.QueryRow(selectSQL, username)
	err := row.Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("ユーザーデータが獲得できませんでした: %w", err)
	}

	return userID, nil
}
