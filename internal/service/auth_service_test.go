package service

import (
	"errors"
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

func newTestAuthService(repo AuthRepositoryInterface) *AuthService {
	jwtService := NewJWTService("test-secret")
	return NewAuthService(repo, jwtService)
}

func newTestUser(t *testing.T, username, password string) *model.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		t.Fatal(err)
	}

	return &model.User{
		ID:           1,
		Username:     username,
		PasswordHash: string(hash),
		Role:         "user",
	}
}

const (
	username string = "mario"
	password string = "123456"
)

func TestLogin_Success(t *testing.T) {
	user := newTestUser(t, username, password)

	repo := &fakeAuthRepository{
		user: user,
	}

	service := newTestAuthService(repo)

	token, err := service.Login(username, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &fakeAuthRepository{
		user: nil,
		err:  errors.New("user not found"),
	}

	service := newTestAuthService(repo)

	token, err := service.Login(username, password)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if token != "" {
		t.Errorf("expected empty, got %s", token)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	user := newTestUser(t, username, password)

	repo := &fakeAuthRepository{
		user: user,
	}

	service := newTestAuthService(repo)

	wrongPassword := "123123"
	token, err := service.Login(username, wrongPassword)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}
}
