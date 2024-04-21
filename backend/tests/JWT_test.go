package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	userID := 123
	login := "testuser"

	token, err := models.GenerateJWT(userID, login)
	if err != nil {
		t.Fatalf("GenerateJWT() failed: %v", err)
	}

	if token == "" {
		t.Error("GenerateJWT() failed: empty token")
	}
}

func TestParseJWT(t *testing.T) {
	userID := 123
	login := "testuser"
	token, err := models.GenerateJWT(userID, login)
	if err != nil {
		t.Fatalf("GenerateJWT() failed: %v", err)
	}

	parsedUserID, parsedLogin, err := models.ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT() failed: %v", err)
	}

	if parsedUserID != float64(userID) {
		t.Errorf("ParseJWT() failed: parsed user ID mismatch. Expected: %d, Got: %f", userID, parsedUserID)
	}

	if parsedLogin != login {
		t.Errorf("ParseJWT() failed: parsed login mismatch. Expected: %s, Got: %s", login, parsedLogin)
	}

	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"name":   login,
		"exp":    time.Now().Add(-1 * time.Hour).Unix(),
	}).SignedString([]byte(models.JWTSecretKey))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	_, _, err = models.ParseJWT(expiredToken)
	if err == nil {
		t.Error("ParseJWT() failed to detect expired token")
	}
}
