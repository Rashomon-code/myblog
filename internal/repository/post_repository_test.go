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
		require.NoError(t, err)
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

func TestGetTitleByUserID(t *testing.T) {
	postRepo := NewPostRepository(testDB)
	authRepo := NewAuthRepository(testDB)

	t.Run("pagination and filtering scenarios", func(t *testing.T) {
		resetTable(t)
		defer resetTable(t)

		authRepo.CreateUserWithProfile("User1", "pass")
		authRepo.CreateUserWithProfile("User2", "pass")
		postRepo.CreatePost(1, "User1 - Post 1", "")
		postRepo.CreatePost(1, "User1 - Post 2", "")
		postRepo.CreatePost(1, "User1 - Post 3", "")
		postRepo.CreatePost(2, "User2 - Post 1", "")

		tests := []struct {
			name           string
			userID         int64
			page           int
			pageSize       int
			wantPageCount  int
			wantTotal      int64
			wantFirstTitle string
		}{
			{
				name:           "User1 Page 1",
				userID:         1,
				page:           1,
				pageSize:       2,
				wantPageCount:  2,
				wantTotal:      3,
				wantFirstTitle: "User1 - Post 3",
			},
			{
				name:           "User1 Page 2",
				userID:         1,
				page:           2,
				pageSize:       2,
				wantPageCount:  1,
				wantTotal:      3,
				wantFirstTitle: "User1 - Post 1",
			},
			{
				name:          "invalid user",
				userID:        10,
				page:          1,
				pageSize:      10,
				wantPageCount: 0,
				wantTotal:     0,
			},
			{
				name:          "invalid page",
				userID:        2,
				page:          2,
				pageSize:      2,
				wantPageCount: 0,
				wantTotal:     1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				posts, total, err := postRepo.GetTitleByUserID(tt.userID, tt.page, tt.pageSize)

				require.NoError(t, err)
				assert.Equal(t, tt.wantTotal, total)
				assert.Len(t, posts, tt.wantPageCount)

				if tt.wantPageCount > 0 {
					assert.Equal(t, tt.wantFirstTitle, posts[0].Title)
				}
			})
		}
	})
}
