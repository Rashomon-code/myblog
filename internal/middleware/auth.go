package middleware

import (
	"net/http"
	"strings"

	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	jwtService *service.JWTService
}

func NewMiddleware(jwtService *service.JWTService) *Middleware {
	return &Middleware{jwtService: jwtService}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "Authorization が見つかりませんでした。"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "求められる Token ではありません。"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(m.jwtService.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "Token が無効でした。"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "求められる Token ではありません。"})
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ユーザーが見つかりませんでした。"})
			c.Abort()
			return
		}

		userid, ok := claims["id"].(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ログイン中にエラーが起きました。"})
			c.Abort()
			return
		}

		c.Set("userID", userid)
		c.Set("username", username)
		c.Next()
	}
}
