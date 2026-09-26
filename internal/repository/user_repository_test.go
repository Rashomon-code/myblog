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
