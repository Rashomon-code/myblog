package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupFailedTestRouter(mw gin.HandlerFunc) (*gin.Engine, *bool) {
	handlerReached := false

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw)

	r.GET("/protected", func(c *gin.Context) {
		handlerReached = true
	})

	return r, &handlerReached
}

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
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)

	r, handlerReached := setupFailedTestRouter(mw.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil) // "/test" 定義されていない Path にアクセスしようとしても、Middleware に経由します
	// req.Header.Set("", "")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, *handlerReached)
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	r, handlerReached := setupFailedTestRouter(mw.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer this is an expected 401 test header")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, *handlerReached)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)

	wrongJWTService := service.NewJWTService("wrong")
	tokenString, _ := wrongJWTService.GenerateToken("wrong", int64(100), "user")

	r, handlerReached := setupFailedTestRouter(mw.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, *handlerReached) //assert.Equal(t, handlerReached, false)
	// if handlerReached {
	// 	t.Errorf("expected handlerReached false, got %v", handlerReached)
	// }
	// if w.Code != http.StatusUnauthorized {
	// 	t.Errorf("expected status %v, got %d", http.StatusUnauthorized, w.Code)
	// }
	//どちらも内容を比較し、結果を t に記録します
}

func TestRequireRole_Success(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "admin")
	})
	r.Use(mw.RequireRole("admin"))
	r.GET("/role", func(c *gin.Context) {
		handlerReached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/role", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, handlerReached)
}

func TestRequireRole_Forbidden(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
	})
	r.Use(mw.RequireRole("admin"))
	r.GET("/role", func(c *gin.Context) {
		handlerReached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/role", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, handlerReached)
}

func TestRequireRole_MissingRoleContext(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.RequireRole("admin"))
	r.GET("/role", func(c *gin.Context) {
		handlerReached = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/role", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.False(t, handlerReached)
}

func TestOptionalAuthMiddleware_ValidToken(t *testing.T) {
	handlerReached := false
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.OptionalAuthMiddleware())
	r.GET("/optional", func(c *gin.Context) {
		handlerReached = true
		userID, exists := c.Get("userID")
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

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, handlerReached)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptionalAuthMiddleware_NoHeader(t *testing.T) {
	handlerReached := false
	var userID int64
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.OptionalAuthMiddleware())
	r.GET("/optional", func(c *gin.Context) {
		handlerReached = true
		userIDAny, exists := c.Get("userID")
		if exists {
			userID = userIDAny.(int64)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, handlerReached)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), userID)
}

func TestOptionalAuthMiddleware_InvalidToken(t *testing.T) {
	handlerReached := false
	var userID int64
	jwtService := service.NewJWTService("test")
	mw := NewMiddleware(jwtService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw.OptionalAuthMiddleware())
	r.GET("/optional", func(c *gin.Context) {
		handlerReached = true
		userIDAny, exists := c.Get("userID")
		if exists {
			userID = userIDAny.(int64)
		}
		c.Status(http.StatusOK)
	})

	tokenString := "this.is.a.test.token" //期限切れtokenを作成してもいい、かりのclaimsから書く必要があります
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, handlerReached)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), userID)
}
