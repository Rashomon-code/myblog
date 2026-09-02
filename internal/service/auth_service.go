package service

import (
	"errors"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryInterface interface {
	CreateUserWithProfile(username, passwordHash string) error
	GetUserByUsername(username string) (*model.User, error)
}

type AuthService struct {
	repo AuthRepositoryInterface
	jwt  *JWTService
}

func NewAuthService(repo AuthRepositoryInterface, jwt *JWTService) *AuthService {
	return &AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Register(username, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4) //学習のため、最低レベルを使用します
	if err != nil {
		return fmt.Errorf("登録できませんでした: %w", err)
	}

	err = s.repo.CreateUserWithProfile(username, string(passwordHash))
	return err
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("入力に誤りがございます。")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("入力に誤りがございます。")
	}

	token, err := s.jwt.GenerateToken(username, user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
