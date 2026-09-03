package service

import (
	"testing"
	"time"

	"github.com/Rashomon-code/myblog/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_Success(t *testing.T) {
	jwtService := NewJWTService("test")

	tokenString, err := jwtService.GenerateToken("user", int64(100), "admin")
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	if tokenString == "" {
		t.Error("expected non-empty, got empty")
	}
}

func TestParseToken_Success(t *testing.T) {
	jwtService := NewJWTService("test")
	tokenString, err := jwtService.GenerateToken("user", int64(100), "admin")
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	parsedClaims, err := jwtService.ParseToken(tokenString)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	if parsedClaims.Username != "user" {
		t.Errorf("expected Username %q, got %q", "user", parsedClaims.Username)
	}
	if parsedClaims.UserID != 100 {
		t.Errorf("expected UserID %d, got %d", 100, parsedClaims.UserID)
	}
	if parsedClaims.Role != "admin" {
		t.Errorf("expected Username %q, got %q", "admin", parsedClaims.Role)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	jwtService := NewJWTService("correct")
	tokenString, _ := jwtService.GenerateToken("user", int64(100), "admin")

	wrongService := NewJWTService("wrong")
	_, err := wrongService.ParseToken(tokenString)
	if err == nil {
		t.Error("expected err, got no err")
	}
}

func TestParseToken_Expired(t *testing.T) {
	expiredClaims := model.Claims{
		UserID:   int64(100),
		Username: "user",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	jwtService := NewJWTService("test")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, _ := token.SignedString([]byte("test"))

	_, err := jwtService.ParseToken(tokenString)
	if err == nil {
		t.Error("expected err, got no err")
	}
}
