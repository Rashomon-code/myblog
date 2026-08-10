package service

import (
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
