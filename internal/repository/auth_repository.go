package repository

import (
	"database/sql"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUserWithProfile(username, passwordHash string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertSQL := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var userID int64
	err = tx.QueryRow(insertSQL, username, passwordHash).Scan(&userID)
	if err != nil {
		return fmt.Errorf("登録にエラーが起きました: %w", err)
	}

	insertProfileSQL := `
		INSERT INTO user_profiles (user_id, display_name, bio)
		VALUES ($1, $2, $3)
	`
	defaultName := fmt.Sprintf("ユーザー %d", userID)
	_, err = tx.Exec(insertProfileSQL, userID, defaultName, "")
	if err != nil {
		return fmt.Errorf("プロフィール設定できませんでした: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("登録完了しました。")
	return nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	selectSQL := `SELECT id, username, password_hash, role FROM users WHERE username = $1`

	err := r.db.QueryRow(selectSQL, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
