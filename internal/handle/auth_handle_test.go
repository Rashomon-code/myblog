package handle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

type mockAuthService struct {
	AuthService

	registerFn func(username string, password string) error
	loginFn    func(username, password string) (string, error)
}

func (m *mockAuthService) Register(username string, password string) error {
	return m.registerFn(username, password)
}

func (m *mockAuthService) Login(username, password string) (string, error) {
	return m.loginFn(username, password)
}

func setupTestRouter(method, url string, handler gin.HandlerFunc, middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	for _, m := range middlewares {
		if m != nil {
			r.Use(m)
		}
	}

	r.Handle(method, url, handler)
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
			r := setupTestRouter(http.MethodPost, "/test", h.Register)

			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("expected status code %d, got %d, %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		mockLoginFn func(username, password string) (string, error)
		wantCode    int
	}{
		{
			name:    "success",
			reqBody: `{"username": "testname", "password": "123456"}`,
			mockLoginFn: func(username, password string) (string, error) {
				return "", nil
			},
			wantCode: http.StatusOK,
		},
		{
			name:    "invalid json body",
			reqBody: `{"username": "testname", "password": }`,
			mockLoginFn: func(username, password string) (string, error) {
				return "", nil
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "invalid user",
			reqBody: `{"username": "testname", "password": "654321"}`,
			mockLoginFn: func(username, password string) (string, error) {
				return "", service.ErrLogin
			},
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockAuthService{loginFn: tt.mockLoginFn}
			h := NewAuthHandle(mockService)
			r := setupTestRouter(http.MethodPost, "/test", h.Login)
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			if w.Code != tt.wantCode {
				t.Errorf("expected status code %d, got %d, %s", tt.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
