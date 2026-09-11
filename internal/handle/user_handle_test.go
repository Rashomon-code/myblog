package handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	UserService

	getProfilefn func(userID int64) (*model.UserProfile, error)
	updateRolefn func(operatorID, userID int64, newRole string) error
}

func (m *mockUserService) GetProfile(userID int64) (*model.UserProfile, error) {
	return m.getProfilefn(userID)
}

func (m *mockUserService) UpdateRole(operatorID, userID int64, newRole string) error {
	return m.updateRolefn(operatorID, userID, newRole)
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

func TestUpdateRole(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		wantCode     int
		reqBody      string
		middleware   func(*gin.Context)
		updateRolefn func(operatorID, userID int64, newRole string) error
	}{
		{
			name:     "success",
			paramID:  "100",
			wantCode: http.StatusOK,
			reqBody:  `{"role": "user"}`,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(200))
				c.Next()
			},
			updateRolefn: func(operatorID, userID int64, newRole string) error {
				assert.Equal(t, int64(200), operatorID)
				assert.Equal(t, int64(100), userID)
				assert.Equal(t, "user", newRole)
				return nil
			},
		},
		{
			name:     "invalid user ID",
			paramID:  "abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing user ID",
			paramID:  "100",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:    "invalid JSON",
			paramID: "100",
			reqBody: `{"": ""}`,
			middleware: func(ctx *gin.Context) {
				ctx.Set("userID", int64(200))
				ctx.Next()
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "service error",
			paramID:  "100",
			wantCode: http.StatusBadRequest,
			reqBody:  `{"role": "user"}`,
			middleware: func(c *gin.Context) {
				c.Set("userID", int64(200))
				c.Next()
			},
			updateRolefn: func(operatorID, userID int64, newRole string) error {
				return errors.New("permission denied")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{
				updateRolefn: tt.updateRolefn,
			}
			h := NewUserHandle(mockSvc)
			r := setupTestRouter(http.MethodPut, "/test/:id", h.UpdateRole, tt.middleware)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/test/"+tt.paramID, strings.NewReader(tt.reqBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
