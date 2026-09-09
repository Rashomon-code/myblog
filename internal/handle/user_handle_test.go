package handle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/gin-gonic/gin"
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
		setupContext func(c *gin.Context)
		getProfilefn func(userID int64) (*model.UserProfile, error)
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			setupContext: func(c *gin.Context) {
				c.Set("userID", int64(100))
			},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				if userID != 100 {
					return nil, fmt.Errorf("unexpected userID: %d", userID)
				}
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:         "no token",
			wantCode:     http.StatusUnauthorized,
			setupContext: func(c *gin.Context) {},
			getProfilefn: func(userID int64) (*model.UserProfile, error) {
				return &model.UserProfile{UserID: userID}, nil
			},
		},
		{
			name:     "wrong user",
			wantCode: http.StatusBadRequest,
			setupContext: func(c *gin.Context) {
				c.Set("userID", int64(1000))
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
			gin.SetMode(gin.TestMode)
			r := gin.New()
			m := &mockUserService{getProfilefn: tt.getProfilefn}
			h := NewUserHandle(m)
			r.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})
			r.GET("/test", h.MyPage)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)
			if w.Code != tt.wantCode {
				t.Errorf("expected sataus code %d, got %d, %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
