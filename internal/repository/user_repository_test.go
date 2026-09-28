package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserProfile(t *testing.T) {
	resetTable(t)
	defer resetTable(t)

	err := initAdmin(testDB)
	require.NoError(t, err)

	userRepo := NewUserRepository(testDB)

	t.Run("no displayname", func(t *testing.T) {
		userProfile, err := userRepo.GetUserProfile(1)
		require.NoError(t, err)
		assert.Equal(t, "ユーザー 1", userProfile.DisplayName)
	})

	t.Run("success", func(t *testing.T) {
		userRepo.UpdateUserProfile(1, "Admin", "")
		userProfile, err := userRepo.GetUserProfile(1)
		require.NoError(t, err)
		assert.Equal(t, "Admin", userProfile.DisplayName)
		assert.Equal(t, "", userProfile.Bio)
	})

	t.Run("no profile", func(t *testing.T) {
		userProfile, err := userRepo.GetUserProfile(2)
		require.NoError(t, err)
		assert.Equal(t, "ユーザー 2", userProfile.DisplayName)
		assert.Equal(t, "まだ何もありません", userProfile.Bio)
	})
}

func TestUpdateRole(t *testing.T) {
	resetTable(t)
	defer resetTable(t)

	authRepo := NewAuthRepository(testDB)
	userRepo := NewUserRepository(testDB)
	err := authRepo.CreateUserWithProfile("testuser", "hashedpass")
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:   "success",
			userID: 1,
		},
		{
			name:    "invalid user",
			userID:  10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.UpdateRole(tt.userID, "admin")
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {

}
