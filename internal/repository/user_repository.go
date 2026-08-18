package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserProfile(userID int64) (*model.UserProfile, error) {
	selectSQL := `SELECT display_name, bio FROM user_profiles WHERE user_id = $1`

	displayName := fmt.Sprintf("ユーザー %d", userID)
	bio := "まだ何もありません"

	row := r.db.QueryRow(selectSQL, userID)
	err := row.Scan(&displayName, &bio)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			return nil, err
		}
	}

	userProfile := model.UserProfile{
		UserID:      userID,
		DisplayName: displayName,
		Bio:         bio,
	}

	return &userProfile, nil
}

func (r *UserRepository) UpdateRole(userID int64, newRole string) error {
	updateSQL := `UPDATE users SET role = $1 WHERE id = $2`
	result, err := r.db.Exec(updateSQL, newRole, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("ユーザーが見つかりませんでした。")
	}

	return nil
}
