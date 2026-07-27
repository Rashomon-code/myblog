package handle

import (
	"net/http"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/Rashomon-code/myblog/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandle struct {
	authService *service.AuthService
}

func NewAuthHandle(s *service.AuthService) *AuthHandle {
	return &AuthHandle{authService: s}
}

func (a *AuthHandle) RegisterAPI(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": err.Error()})
		return
	}

	err := a.authService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登録しました"})
}

func (a *AuthHandle) LoginAPI(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": err.Error()})
		return
	}

	err := a.authService.Login(req.Username, req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ログインしました"})
}
