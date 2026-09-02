package service

import (
	"testing"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_Success(t *testing.T) {
	jwtService := NewJWTService("test")

	tokenString, err := jwtService.GenerateToken("user", int64(100), "admin")
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	parsedClaims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, parsedClaims, func(t *jwt.Token) (any, error) {
		return []byte("test"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("expected no err, got %v", err)
	}

	if parsedClaims.UserID != 100 {
		t.Errorf("expected UserID 100, got %d", parsedClaims.UserID)
	}

	if parsedClaims.Username != "user" {
		t.Errorf("expected Username %q, got %q", "user", parsedClaims.Username)
	}

	if parsedClaims.Role != "admin" {
		t.Errorf("expected Role %q, got %q", "admin", parsedClaims.Role)
	}
}

func TestGenerateToken_InvalidSecret(t *testing.T) {
	jwtService := NewJWTService("correct-secret")
	tokenString, _ := jwtService.GenerateToken("user", 100, "admin")

	parsedClaims := &model.Claims{}
	_, err := jwt.ParseWithClaims(tokenString, parsedClaims, func(t *jwt.Token) (any, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Error("expected err,got no err")
	}
}
