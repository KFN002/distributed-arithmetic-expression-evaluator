package models

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const JWTSecretKey = "kogda_ya_na_pochte_sluzhil_yamshikom_ko_mne_postuchalsya_kosmatiy_geolog"

func GenerateJWT(userID int, login string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"name":   login,
		"exp":    now.Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSecretKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["userID"].(int)
	if !ok {
		return 0, "", fmt.Errorf("user ID not found in token claims")
	}
	name, ok := claims["name"].(string)
	if !ok {
		return 0, "", fmt.Errorf("name not found in token claims")
	}

	return userID, name, nil
}
