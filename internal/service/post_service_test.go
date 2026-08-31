package service

import (
	"errors"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
)

//ムズイ！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！

type mockPostRepository struct {
	PostRepository

	GetPostDetailFunc func(postID int64) (model.PostDetail, error)
	DeletePostFunc    func(postID int64) error
}

func (m *mockPostRepository) DeletePost(postID int64) error {
	return m.DeletePostFunc(postID)
}

func (m *mockPostRepository) GetPostDetail(postID int64) (model.PostDetail, error) {
	return m.GetPostDetailFunc(postID)
}

func TestCreatePost(t *testing.T) {
	s := NewPostService(nil)

	testCases := []string{
		"",
		" ",
		"\n",
	}

	for _, tc := range testCases {
		err := s.CreatePost(1, tc, "テスト")
		if err == nil {
			t.Fatalf("タイトル [%s]、エラーのはずですが、結果は nil でした", tc)
		}

		expectedMsg := "タイトルが入力されていません"
		if err.Error() != expectedMsg {
			t.Errorf("エラーメッセージが相違します: [%s]のはずですが、 [%s]でした", expectedMsg, err.Error())
		}
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
