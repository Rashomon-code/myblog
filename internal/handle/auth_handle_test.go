package handle

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	AuthService

	registerFn func(username string, password string) error
}

func (m *mockAuthService) Register(username string, password string) error {
	return m.registerFn(username, password)
}

func setupTestRouter(handler gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/test", handler)
	return r
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        string
		mockRegisterFn func(username, password string) error
		wantCode       int
	}{
		{
			name:    "success",
			reqBody: `{"username": "test", "password": "123"}`,
			mockRegisterFn: func(username, password string) error {
				return nil
			},
			wantCode: http.StatusOK,
		},
		{
			name:    "invalid json body",
			reqBody: `{"username": , "password": }`,
			mockRegisterFn: func(username, password string) error {
				return nil
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "service error",
			reqBody: `{"username": "test", "password": "123"}`,
			mockRegisterFn: func(username, password string) error {
				return errors.New("db error")
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAuthService{registerFn: tt.mockRegisterFn}
			h := NewAuthHandle(mockService)
			r := setupTestRouter(h.Register)

			req, _ := http.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("expected status code %d, got %d, %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
