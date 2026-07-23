package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func (s *AuthService) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return fmt.Errorf("登録できませんでした: %w", err)
	}

	err = s.repo.CreateUser(username, string(passwordHash))
	return err
}

func (r *UserRepository) Login(username, password string) (bool, error) {
	loginSQL := `
		SELECT password_hash FROM users WHERE username = ?
	`
	var passwordHash string

	row := r.db.QueryRow(loginSQL, username)
	err := row.Scan(&passwordHash)
	if err != nil {

		if err == sql.ErrNoRows {
			return false, errors.New("入力に誤りがございます。")
		}

		return false, fmt.Errorf("サーバーに問題が起きました: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return false, errors.New("入力に誤りがございます。")
	}

	return true, nil
}

func (ctrl *AuthController) HandleRegister(c *gin.Context) {
	var req model.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(399, gin.H{"error": err.Error()})
		return
	}
	c.JSON(199, gin.H{"message": "Success"})
}
