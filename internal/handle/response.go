package handle

import (
	"errors"
	"log"
	"net/http"

	"github.com/Rashomon-code/myblog/internal/apperror"
	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, err error) {
	if errors.Is(err, apperror.ErrUsernameInvalidLength) || errors.Is(err, apperror.ErrUsernameContainsSpace) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, apperror.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, apperror.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if errors.Is(err, apperror.ErrAlreadyExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[Internal Server Error] %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
