package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/Rashomon-code/myblog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameInvalidLength = errors.New("username must be between 3 and 20 characters")
	ErrUsernameContainsSpace = errors.New("username cannot contain spaces")
	ErrDatabase              = errors.New("database error")
	ErrLogin                 = errors.New("入力に誤りがございます")
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

// 正規表現で制限する方法
// var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
// func validateUsernameFormat(username string) error {
// 	if !usernameRegex.MatchString(username) {
// 		return errors.New("Format error")
// 	}
// 	return nil
// }

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return ErrUsernameInvalidLength
	}

	if strings.IndexFunc(username, unicode.IsSpace) != -1 {
		return ErrUsernameContainsSpace
	}

	return nil
}

func (s *AuthService) Register(username, password string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

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
		return "", ErrLogin
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrLogin
	}

	token, err := s.jwt.GenerateToken(username, user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
