package handle

import (
	"net/http"

	"github.com/Rashomon-code/myblog/internal/apperror"
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
		RespondWithError(c, apperror.ErrInvalidInput)
		return
	}

	err := a.authService.Register(req.Username, req.Password)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (a *AuthHandle) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, apperror.ErrInvalidInput)
		return
	}

	token, err := a.authService.Login(req.Username, req.Password)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successfully", "token": token})
}
