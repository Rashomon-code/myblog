package service

import (
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthRepository struct {
	user *model.User
	err  error
}

func (f *fakeAuthRepository) GetUserByUsername(username string) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeAuthRepository) CreateUserWithProfile(username, passwordHash string) error {
	return f.err
}

func TestLogin(t *testing.T) {
	username := "mario"
	password := "123456"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		t.Fatal(err)
	}

	repo := &fakeAuthRepository{
		user: &model.User{
			ID:           1,
			Username:     username,
			PasswordHash: string(hash),
			Role:         "user",
		},
	}

	jwtService := NewJWTService("test-secret")

	service := NewAuthService(repo, jwtService)

	token, err := service.Login(username, password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}
