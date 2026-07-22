package repository

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return fmt.Errorf("創建失敗: %w", err)
	}

	insertSQL := `
		INSERT INTO users (username, password_hash)
		VALUES (?, ?)
	`

	_, err = r.db.Exec(insertSQL, username, string(passwordHash))
	if err != nil {
		return fmt.Errorf("創建失敗: %w", err)
	}

	fmt.Println(username, "創建成功")
	return nil
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
