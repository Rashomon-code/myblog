package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Success(t *testing.T) {
	handlerReached := false

	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		handlerReached = true
		userID, exists := c.Get("userID")
		//userID.(int64) != 100 ok, userID != int64(100) ok, but userID != 100 fail
		if !exists || userID != int64(100) {
			c.Status(http.StatusInternalServerError)
			return
		}

		username, exists := c.Get("username")
		if !exists || username != "testuser" {
			c.Status(http.StatusInternalServerError)
			return
		}

		role, exists := c.Get("role")
		if !exists || role != "user" {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	tokenString, _ := jwtService.GenerateToken("testuser", int64(100), "user")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, handlerReached)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.AuthMiddleware())
	r.GET("/protected", func(ctx *gin.Context) {
		handlerReached = true
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("", "")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if handlerReached == true {
		t.Errorf("expected handlerReached false, got %v", handlerReached)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.AuthMiddleware())
	r.GET("/protected", func(ctx *gin.Context) {
		handlerReached = true
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer this is an expected 401 test header")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if handlerReached {
		t.Errorf("expected handlerReached false, got %v", handlerReached)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)

	wrongJWTService := service.NewJWTService("wrong")
	tokenString, _ := wrongJWTService.GenerateToken("wrong", int64(100), "user")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.AuthMiddleware())
	r.GET("/protected", func(ctx *gin.Context) {
		handlerReached = true
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, handlerReached) //assert.Equal(t, handlerReached, false)
	// if handlerReached {
	// 	t.Errorf("expected handlerReached false, got %v", handlerReached)
	// }
	// if w.Code != http.StatusUnauthorized {
	// 	t.Errorf("expected status %v, got %d", http.StatusUnauthorized, w.Code)
	// }
	//どちらも内容を比較し、結果を t に記録します
}
