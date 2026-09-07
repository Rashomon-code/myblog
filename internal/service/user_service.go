package service

import (
	"errors"

	"github.com/Rashomon-code/myblog/internal/model"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	GetUserProfile(userID int64) (*model.UserProfile, error)
	UpdateRole(userID int64, newRole string) error
	GetAllUsers() ([]model.UserResponse, error)
	UpdateUserProfile(userID int64, displayName, bio string) error
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfileService(userID int64) (*model.UserProfile, error) {
	return s.repo.GetUserProfile(userID)
}

func (s *UserService) UpdateRoleService(operatorID, userID int64, newRole string) error {
	if newRole != "admin" && newRole != "user" {
		return errors.New("無効なタイプ")
	}

	if operatorID == userID {
		return errors.New("自分の権限を変更することができません")
	}

	return s.repo.UpdateRole(userID, newRole)
}

func (s *UserService) GetAllUsersService() ([]model.UserResponse, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateProfileService(userID int64, displayName, bio string) error {
	return s.repo.UpdateUserProfile(userID, displayName, bio)
}
