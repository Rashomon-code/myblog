package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/Rashomon-code/myblog/internal/apperror"
	"github.com/Rashomon-code/myblog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryInterface interface {
	CreateUserWithProfile(username, passwordHash string) error
	GetUserByUsername(username string) (*model.User, error)
	SaveRefreshToken(userID int64, refreshToken string, expiresAt time.Time) error
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
		return apperror.ErrUsernameInvalidLength
	}

	if strings.IndexFunc(username, unicode.IsSpace) != -1 {
		return apperror.ErrUsernameContainsSpace
	}

	return nil
}

func (s *AuthService) Register(username, password string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost) //学習のため、最低レベルを使用します
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.repo.CreateUserWithProfile(username, string(passwordHash))
	return err
}

func (s *AuthService) Login(username, password string) (*model.TokenPair, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	accessToken, err := s.jwt.GenerateToken(username, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	err = s.repo.SaveRefreshToken(user.ID, refreshToken, expiresAt)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
