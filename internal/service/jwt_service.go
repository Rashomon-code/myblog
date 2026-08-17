package service

import (
	"time"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	Secret string
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{Secret: secret}
}

func (j *JWTService) GenerateToken(username string, userID int64, role string) (string, error) {
	claims := model.Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
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
