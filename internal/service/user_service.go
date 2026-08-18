package service

import (
	"errors"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
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
