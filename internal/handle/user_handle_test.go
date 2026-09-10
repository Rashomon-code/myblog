package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	UserService

	getProfilefn func(userID int64) (*model.UserProfile, error)
}

func (m *mockUserService) GetProfile(userID int64) (*model.UserProfile, error) {
	return m.getProfilefn(userID)
}

func TestMypage(t *testing.T) {
	tests := []struct {
		name         string
		wantCode     int
		middleware   func(c *gin.Context)
		getProfilefn func(userID int64) (*model.UserProfile, error)
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				if userID != 100 {
					return nil, fmt.Errorf("unexpected userID: %d", userID)
				}
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "no token",
			wantCode: http.StatusUnauthorized,
			middleware: func(c *gin.Context) {
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "wrong user",
			wantCode: http.StatusBadRequest,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(1000))
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				if userID != 100 {
					return nil, fmt.Errorf("unexpected userID: %d", userID)
				}
				return &model.UserProfile{UserID: userID}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockUserService{getProfilefn: tt.getProfilefn}
			h := NewUserHandle(m)

			r := setupTestRouter(http.MethodGet, "/test", h.MyPage, tt.middleware)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)
			if w.Code != tt.wantCode {
				t.Errorf("expected sataus code %d, got %d, %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		paramID      string
		wantCode     int
		wantIsMe     bool
		middleware   func(c *gin.Context)
		getProfilefn func(userID int64) (*model.UserProfile, error)
	}{
		{
			name:     "success self",
			paramID:  "100",
			wantCode: http.StatusOK,
			wantIsMe: true,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "success orter",
			paramID:  "100",
			wantCode: http.StatusOK,
			wantIsMe: false,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(200))
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "invalid URL ID format",
			paramID:  "abc",
			wantCode: http.StatusUnauthorized,
			wantIsMe: false,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(100))
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "success view",
			paramID:  "100",
			wantCode: http.StatusOK,
			wantIsMe: false,
			middleware: func(c *gin.Context) {
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "user not found",
			paramID:  "100",
			wantCode: http.StatusBadRequest,
			wantIsMe: false,
			middleware: func(c *gin.Context) {
				c.Next()
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return nil, errors.New("user not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{getProfilefn: tt.getProfilefn}
			h := NewUserHandle(mockSvc)

			r := setupTestRouter(http.MethodGet, "/test/:id", h.GetUserProfile, tt.middleware)

			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.paramID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantCode, w.Code, w.Body.String())
			if tt.wantCode == http.StatusOK {
				var res model.UserProfileResponse
				err := json.Unmarshal(w.Body.Bytes(), &res)
				require.NoError(t, err)
				assert.Equal(t, tt.wantIsMe, res.IsMe)
			}
		})
	}
}
