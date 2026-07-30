package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{Secret: secret}
}

func (j *JWTService) GenerateToken(username string, userID int64) (string, error) {
	claims := jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp": time.Now().Add(
			time.Hour * 24,
		).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(j.Secret),
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
