package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Rashomon-code/myblog/internal/apperror"
	"github.com/Rashomon-code/myblog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepository struct {
	AuthRepositoryInterface

	createUserFn func(username, passwordHash string) error
	getUserFn    func(username string) (*model.User, error)
	saveTokenFn  func(userID int64, refreshToken string, expiresAt time.Time) error
}

func (r *mockAuthRepository) CreateUserWithProfile(username, passwordHash string) error {
	return r.createUserFn(username, passwordHash)
}

func (r *mockAuthRepository) GetUserByUsername(username string) (*model.User, error) {
	return r.getUserFn(username)
}

func (r *mockAuthRepository) SaveRefreshToken(userID int64, refreshToken string, expiresAt time.Time) error {
	return r.saveTokenFn(userID, refreshToken, expiresAt)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name             string
		username         string
		password         string
		mockCreateUserFn func(username, passwordHash string) error
		wantErr          error
	}{
		{
			name:     "success",
			username: "testuser",
			password: "123456",
			mockCreateUserFn: func(username, passwordHash string) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name:             "invalid username space",
			username:         " wrong ",
			password:         "123456",
			mockCreateUserFn: nil,
			wantErr:          apperror.ErrUsernameContainsSpace,
		},
		{
			name:             "invalid username",
			username:         "wr",
			password:         "123456",
			mockCreateUserFn: nil,
			wantErr:          apperror.ErrUsernameInvalidLength,
		},
		{
			name:     "repository error",
			username: "testuser",
			password: "123456",
			mockCreateUserFn: func(username, passwordHash string) error {
				return apperror.ErrDatabase
			},
			wantErr: apperror.ErrDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAuthRepository{
				createUserFn: tt.mockCreateUserFn,
			}
			jwtService := NewJWTService("test")
			mockS := NewAuthService(r, jwtService)

			err := mockS.Register(tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name            string
		username        string
		password        string
		mockGetUserfn   func(username string) (*model.User, error)
		mockSaveTokenFn func(userID int64, refreshToken string, expiresAt time.Time) error
		wantErr         error
	}{
		{
			name:     "success",
			username: "testuser",
			password: "123456",
			mockGetUserfn: func(username string) (*model.User, error) {
				hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 4)
				if err != nil {
					t.Fatal(err)
				}
				return &model.User{
					ID:           100,
					Username:     "testuser",
					PasswordHash: string(hash),
					Role:         "admin",
				}, nil
			},
			mockSaveTokenFn: func(userID int64, refreshToken string, expiresAt time.Time) error {
				if userID != 100 {
					return errors.New("failed to save token")
				}
				return nil
			},
			wantErr: nil,
		},
		{
			name: "wrong user",
			mockGetUserfn: func(username string) (*model.User, error) {
				return nil, apperror.ErrInvalidCredentials
			},
			wantErr: apperror.ErrInvalidCredentials,
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "123456",
			mockGetUserfn: func(username string) (*model.User, error) {
				hash, err := bcrypt.GenerateFromPassword([]byte("654321"), 4)
				if err != nil {
					t.Fatal(err)
				}
				return &model.User{
					ID:           100,
					Username:     "testuser",
					PasswordHash: string(hash),
					Role:         "admin",
				}, apperror.ErrInvalidCredentials
			},
			wantErr: apperror.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mockAuthRepository{
				getUserFn:   tt.mockGetUserfn,
				saveTokenFn: tt.mockSaveTokenFn,
			}
			jwtService := NewJWTService("test")
			mockS := NewAuthService(r, jwtService)

			_, err := mockS.Login(tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
