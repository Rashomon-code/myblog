package middleware

import (
	"net/http"
	"strings"

	"github.com/Rashomon-code/myblog/internal/model"
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
		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(t *jwt.Token) (any, error) {
			return []byte(m.jwtService.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "無効な Token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*model.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "求められる Token ではありません。"})
			c.Abort()
			return
		}

		//戻す際数字はデフォルトのfloat (jwt.MapClaims{} を使用する場合)
		// userid, ok := claims["id"].(float64)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ログイン中にエラーが起きました。"})
		// 	c.Abort()
		// 	return
		// }
		// username, ok := claims["username"].(string)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"エラー": "ユーザーが見つかりませんでした。"})
		// 	c.Abort()
		// 	return
		// }

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func (m *Middleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"エラー": "認証情報が見つかりません。"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"エラー": "アクセス権限がありません。"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 && parts[0] != "Bearer" {
			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(t *jwt.Token) (any, error) {
				return []byte(m.jwtService.Secret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(*model.Claims); ok {
					c.Set("userID", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("role", claims.Role)
				}
			}
		}
		c.Next()
	}
}
