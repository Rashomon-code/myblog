package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	repo := NewPostRepository(testDB)
	t.Run("success", func(t *testing.T) {
		resetTable(t)
		defer resetTable(t)

		err := initAdmin(testDB)
		err = repo.CreatePost(1, "test", "")
		require.NoError(t, err)

		var title string
		err = testDB.QueryRow("SELECT title FROM posts WHERE user_id = 1").Scan(&title)
		require.NoError(t, err)
		assert.Equal(t, "test", title)
	})

	t.Run("invalid user", func(t *testing.T) {
		resetTable(t)
		defer resetTable(t)

		err := repo.CreatePost(1, "test", "")
		require.Error(t, err)
	})
}
