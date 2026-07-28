package service

import (
	"errors"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UserRepository
	jwt  *JWTService
}

func NewAuthService(repo *repository.UserRepository, jwt *JWTService) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return fmt.Errorf("登録できませんでした: %w", err)
	}

	err = s.repo.CreateUser(username, string(passwordHash))
	return err
}

func (s *AuthService) Login(username, password string) (string, error) {
	passwordHash, err := s.repo.GetPasswordHash(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", errors.New("入力に誤りがございます。")
	}

	token, err := s.jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}
