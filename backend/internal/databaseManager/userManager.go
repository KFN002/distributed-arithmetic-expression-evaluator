package databaseManager

import (
	"database/sql"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
)

func SignUpUser(login string, password string) error {
	var count int
	err := DB.DB.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", login).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("User with this nickname already exists")
	}

	result, err := DB.DB.Exec("INSERT INTO users (login, password) VALUES (?, ?)", login, password)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	err = DB.AddOperations(int(userID))
	if err != nil {
		return err
	}

	return nil

}

func LogInUser(login string, password string) (string, error) {
	var storedPassword string
	var userID int

	err := DB.DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login).Scan(&userID, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Incorrect username")
		}
		return "", err
	}

	if storedPassword != password {
		return "", fmt.Errorf("Incorrect password")
	}

	tokenString, err := models.GenerateJWT(userID, login)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}
