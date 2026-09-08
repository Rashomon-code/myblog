package handle

import (
	"log"
	"net/http"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/gin-gonic/gin"
)

type AuthHandle struct {
	authService AuthService
}

func NewAuthHandle(s AuthService) *AuthHandle {
	return &AuthHandle{authService: s}
}

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

func (a *AuthHandle) Register(c *gin.Context) {
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

func (a *AuthHandle) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"エラー": err.Error()})
		return
	}

	token, err := a.authService.Login(req.Username, req.Password)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"エラー": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ログインしました", "token": token})
}
