package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

func (j *JWTService) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
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
		[]byte(j.secret),
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
