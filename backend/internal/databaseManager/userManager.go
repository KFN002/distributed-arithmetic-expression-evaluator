package databaseManager

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const JWTSecretKey = "kogda_ya_na_pochte_sluzhil_yamshikom_ko_mne_postuchalsya_kosmatiy_geolog"

func SignUpUser(login string, password string) error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", login).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("User with this nickname already exists")
	}

	_, err = DB.Exec("INSERT INTO users (login, password) VALUES (?, ?)", login, password)
	if err != nil {
		return err
	}
	return nil
}

func LogInUser(login string, password string) (string, error) {
	var storedPassword string
	var userID int

	err := DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Incorrect username")
		}
		return "", err
	}

	if storedPassword != password {
		return "", fmt.Errorf("Incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}
