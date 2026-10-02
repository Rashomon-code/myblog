package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Rashomon-code/myblog/internal/apperror"
	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUserWithProfile(username, passwordHash string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertSQL := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var userID int64
	err = tx.QueryRow(insertSQL, username, passwordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperror.ErrAlreadyExists
		}

		return fmt.Errorf("failed to insert user: %w", err)
	}

	insertProfileSQL := `
		INSERT INTO user_profiles (user_id, display_name, bio)
		VALUES ($1, $2, $3)
	`
	defaultName := fmt.Sprintf("ユーザー %d", userID)
	_, err = tx.Exec(insertProfileSQL, userID, defaultName, "")
	if err != nil {
		return fmt.Errorf("failed to insert user profile: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	selectSQL := `SELECT id, username, password_hash, role FROM users WHERE username = $1`

	err := r.db.QueryRow(selectSQL, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by username: %w", err)
	}

	return &user, nil
}
