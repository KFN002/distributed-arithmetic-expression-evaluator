package databaseManager

import (
	"database/sql"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// SignUpUser Регистрация пользователя и проверка его данных
func SignUpUser(login string, password string) (error, int) {
	var count int
	err := DB.DB.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", login).Scan(&count)
	if err != nil {
		return err, 0
	}

	if count > 0 {
		return fmt.Errorf("User with this nickname already exists"), 0
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err), 0
	}

	result, err := DB.DB.Exec("INSERT INTO users (login, password) VALUES (?, ?)", login, string(hashedPassword))
	if err != nil {
		return err, 0
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err, 0
	}

	err = DB.AddOperations(int(userID))
	if err != nil {
		return err, 0
	}

	return nil, int(userID)
}

// LogInUser Вход и проверка данных пользователя
func LogInUser(login string, password string) (string, error) {
	var storedPasswordHash string
	var userID int

	err := DB.DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login).Scan(&userID, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Incorrect username")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Incorrect password")
	}

	tokenString, err := models.GenerateJWT(userID, login)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

// CheckUser Проверка на наличие пользователя
func CheckUser(userID float64, userName string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND login = ?)"

	err := DB.DB.QueryRow(query, userID, userName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, nil
}
