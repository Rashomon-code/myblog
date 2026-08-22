package service

import (
	"errors"
	"strings"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/repository"
)

var ErrForbidden = errors.New("権限がありません")

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

func (s *PostService) GetPostTitle(userID int64, page, pageSize int) ([]model.ArticleSummary, int64, error) {
	return s.repo.GetTitleByUserID(userID, page, pageSize)
}

func (s *PostService) PostDetailService(postID int64) (model.PostDetail, error) {
	return s.repo.GetPostDetail(postID)
}

func (s *PostService) DeletePostService(postID, userID int64, userRole string) error {
	post, err := s.repo.GetPostDetail(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID && userRole != "admin" {
		return ErrForbidden
	}

	return s.repo.DeletePost(postID)
}

func (s *PostService) EditPostService(postID int64, title string, content string, userID int64, userRole string) error {
	post, err := s.repo.GetPostDetail(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID && userRole != "admin" {
		return ErrForbidden
	}

	return s.repo.EditPost(postID, title, content)
}

func (s *PostService) GetAllPostsService(page, pageSize int) ([]model.ArticleSummary, int64, error) {
	return s.repo.GetAllPosts(page, pageSize)
}

func (s *PostService) SearchPostService(keyword string) ([]model.ArticleSummary, error) {
	return s.repo.SearchPost(keyword)
}
