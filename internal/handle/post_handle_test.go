package handle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPostService struct {
	PostService

	called       bool
	createPostfn func(userID int64, title, content string) error
	postDetailfn func(postID int64) (model.PostDetail, error)
	deletePostfn func(postID, userID int64, userRole string) error
}

func (m *mockPostService) CreatePost(userID int64, title, content string) error {
	m.called = true
	return m.createPostfn(userID, title, content)
}

func (m *mockPostService) PostDetail(postID int64) (model.PostDetail, error) {
	m.called = true
	return m.postDetailfn(postID)
}

func (m *mockPostService) DeletePost(postID, userID int64, userRole string) error {
	m.called = true
	return m.deletePostfn(postID, userID, userRole)
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name         string
		wantCode     int
		reqBody      string
		wantCalled   bool
		middleware   func(*gin.Context)
		createPostfn func(userID int64, title, content string) error
	}{
		{
			name:       "success",
			wantCode:   http.StatusCreated,
			reqBody:    `{"title": "test post", "content": "test content"}`,
			wantCalled: true,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			createPostfn: func(userID int64, title, content string) error {
				assert.Equal(t, int64(100), userID)
				assert.Equal(t, "test post", title)
				assert.Equal(t, "test content", content)
				return nil
			},
		},
		{
			name:       "missing user",
			wantCode:   http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name:     "missing request",
			wantCode: http.StatusBadRequest,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			wantCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockPostService{
				createPostfn: tt.createPostfn,
			}
			postHandle := NewPostHandle(mockSvc)
			r := setupTestRouter(http.MethodPost, "/test", postHandle.CreatePost, tt.middleware)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.reqBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCalled, mockSvc.called)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestPostDetail(t *testing.T) {
	tests := []struct {
		name         string
		wantCode     int
		paramID      string
		wantCalled   bool
		postDetailfn func(postID int64) (model.PostDetail, error)
	}{
		{
			name:       "success",
			wantCode:   http.StatusOK,
			paramID:    "100",
			wantCalled: true,
			postDetailfn: func(postID int64) (model.PostDetail, error) {
				return model.PostDetail{}, nil
			},
		},
		{
			name:     "bad id",
			wantCode: http.StatusBadRequest,
			paramID:  "a1",
		},
		{
			name:       "service error",
			wantCode:   http.StatusInternalServerError,
			paramID:    "100",
			wantCalled: true,
			postDetailfn: func(postID int64) (model.PostDetail, error) {
				return model.PostDetail{}, errors.New("service err")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockPostService{
				postDetailfn: tt.postDetailfn,
			}
			mockHandle := NewPostHandle(mockSvc)
			r := setupTestRouter(http.MethodGet, "/test/:id", mockHandle.PostDetail)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.paramID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCalled, mockSvc.called)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name         string
		wantCode     int
		paramID      string
		wantCalled   bool
		middleware   func(*gin.Context)
		deletePostfn func(postID, userID int64, userRole string) error
	}{
		{
			name:       "success",
			wantCode:   http.StatusOK,
			paramID:    "100",
			wantCalled: true,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			deletePostfn: func(postID, userID int64, userRole string) error {
				return nil
			},
		},
		{
			name:     "bad post id",
			wantCode: http.StatusBadRequest,
			paramID:  "a1",
		},
		{
			name:       "service error",
			wantCode:   http.StatusInternalServerError,
			paramID:    "100",
			wantCalled: true,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			deletePostfn: func(postID, userID int64, userRole string) error {
				return errors.New("service error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockPostService{
				deletePostfn: tt.deletePostfn,
			}
			mockHandle := NewPostHandle(mockSvc)
			r := setupTestRouter(http.MethodDelete, "/test/:id", mockHandle.DeletePost, tt.middleware)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/test/"+tt.paramID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCalled, mockSvc.called)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
