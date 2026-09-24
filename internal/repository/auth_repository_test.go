package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUserWithProfile(t *testing.T) {
	repo := NewAuthRepository(testDB)

	t.Run("success", func(t *testing.T) {
		resetUsersTable(t)
		defer resetUsersTable(t)

		err := repo.CreateUserWithProfile("testuser", "hashed_pw")
		require.NoError(t, err)

		var userID int64
		err = testDB.QueryRow(`
			SELECT id FROM users WHERE username = 'testuser'	
		`).Scan(&userID)
		require.NoError(t, err)

		var displayName, bio string
		err = testDB.QueryRow(`
			SELECT display_name, bio FROM user_profiles WHERE user_id = $1	
		`, userID).Scan(&displayName, &bio)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("ユーザー %d", userID), displayName)
		assert.Equal(t, "", bio)
	})

	t.Run("false", func(t *testing.T) {
		resetUsersTable(t)
		defer resetUsersTable(t)

		err := repo.CreateUserWithProfile("testuser", "hashed1")
		require.NoError(t, err)

		err = repo.CreateUserWithProfile("testuser", "hashed2")
		require.Error(t, err)

		var count int
		err = testDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		err = testDB.QueryRow("SELECT COUNT(*) FROM user_profiles").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestGetUserByUsername(t *testing.T) {
	repo := NewAuthRepository(testDB)

	t.Run("success", func(t *testing.T) {
		resetUsersTable(t)
		defer resetUsersTable(t)

		err := initAdmin(testDB)
		require.NoError(t, err)

		user, err := repo.GetUserByUsername("admin")
		require.NoError(t, err)

		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "admin", user.Role)

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("admin1234"))
		assert.NoError(t, err)
	})

	t.Run("no user", func(t *testing.T) {
		resetUsersTable(t)
		defer resetUsersTable(t)

		_, err := repo.GetUserByUsername("admin")
		assert.Error(t, err)
	})
}
