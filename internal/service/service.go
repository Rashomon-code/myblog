package service

import (
	"fmt"

	"github.com/Rashomon-code/myblog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return fmt.Errorf("登録できませんでした: %w", err)
	}

	err = s.repo.CreateUser(username, string(passwordHash))
	return err
}
