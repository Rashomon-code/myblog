package service

import (
	"errors"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
)

type mockPostRepository struct {
	PostRepository
	//関数を変数として struct に保存して仕様する前に実作します（関数も変数？関数ポインター？）
	GetPostDetailFunc func(postID int64) (model.PostDetail, error)
	DeletePostFunc    func(postID int64) error
	CreatePostFunc    func(userID int64, title string, content string) error
}

func (m *mockPostRepository) DeletePost(postID int64) error {
	return m.DeletePostFunc(postID)
}

func (m *mockPostRepository) GetPostDetail(postID int64) (model.PostDetail, error) {
	return m.GetPostDetailFunc(postID)
}

func (m *mockPostRepository) CreatePost(userID int64, title string, content string) error {
	return m.CreatePostFunc(userID, title, content)
}

func TestValidateTitle(t *testing.T) {
	testCases := []string{"", " ", "\n", "  \n"}
	for _, tc := range testCases {
		_, err := validateTitle(tc)
		if !errors.Is(err, ErrInvalidTitle) {
			t.Errorf("expected %v, got %v", ErrInvalidTitle, err)
		}
	}
}

func TestCreatePost_InvalidTitle(t *testing.T) {

	service := NewPostService(nil)
	err := service.CreatePost(1, "", "content")
	if !errors.Is(err, ErrInvalidTitle) {
		t.Errorf("expected %v, got %v", ErrInvalidTitle, err)
	}
}

func TestCreatePost_RepoError(t *testing.T) {
	mockRepo := &mockPostRepository{
		CreatePostFunc: func(userID int64, title, content string) error {
			return errors.New("db error")
		},
	}

	service := NewPostService(mockRepo)
	err := service.CreatePost(1, "title", "content")
	if err == nil {
		t.Fatalf("expected error %q, got nil", "db error")
	}
	if err.Error() != "db error" {
		t.Errorf("expected error: %q, got %v", "db error", err)
	}
}

func TestCreatePost_Success(t *testing.T) {
	var capturedTitle string
	mockRepo := &mockPostRepository{
		CreatePostFunc: func(userID int64, title, content string) error {
			capturedTitle = title
			return nil
		},
	}

	service := NewPostService(mockRepo)
	err := service.CreatePost(1, " title   ", "content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedTitle != "title" {
		t.Errorf("expected clean title 'Hello Go', got '%s'", capturedTitle)
	}
}

func TestDeletePost_Success(t *testing.T) {
	mockRepo := &mockPostRepository{
		GetPostDetailFunc: func(postID int64) (model.PostDetail, error) {
			return model.PostDetail{
				ID:     postID,
				UserID: 100,
			}, nil
		},
		DeletePostFunc: func(postID int64) error {
			return nil
		},
	}

	service := NewPostService(mockRepo)
	err := service.DeletePost(1, 100, "user")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDeletePost_Forbidden(t *testing.T) {
	mockRepo := &mockPostRepository{
		GetPostDetailFunc: func(postID int64) (model.PostDetail, error) {
			return model.PostDetail{
				ID:     postID,
				UserID: 100,
			}, nil
		},
	}

	service := NewPostService(mockRepo)

	err := service.DeletePost(1, 200, "user")
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}
