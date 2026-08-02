package service

import (
	"errors"
	"strings"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePostService(userID int64, title string, content string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("タイトルが入力されていません")
	}

	err := s.repo.CreatePost(userID, title, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GetUserMyPage(userID int64) ([]model.ArticleSummary, error) {
	return s.repo.GetTitleByUserID(userID)
}
