package service

import "github.com/Rashomon-code/myblog/internal/repository"

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePostService(userID int, title string, content string) error {
	err := s.repo.CreatePost(userID, title, content)
	return err
}
