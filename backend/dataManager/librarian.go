package dataManager

import (
	"database/sql"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// работа с бд, все просто, везде говорящие имена

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}
}

func GetTimes() ([]int, error) {
	rows, err := DB.Query("SELECT time FROM operations")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, fmt.Errorf("Error querying database: %v", err)
	}
	defer rows.Close()

	var operationTimes []int
	for rows.Next() {
		var time int
		if err := rows.Scan(&time); err != nil {
			log.Println("Error scanning row:", err)
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		operationTimes = append(operationTimes, time)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return operationTimes, nil
}

func CheckDuplicate(expression string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM expressions WHERE expression = ?", expression).Scan(&count)
	if err != nil {
		log.Println("Error querying database:", err)
		return false, fmt.Errorf("Error querying database: %v", err)
	}
	return count > 0, nil
}

func GetId(expression string) (int, error) {
	var id int
	err := DB.QueryRow("SELECT id FROM expressions WHERE expression = ?", expression).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Println("Error querying database:", err)
		return 0, fmt.Errorf("Error querying database: %v", err)
	}

	return id, nil
}

func GetExpressions() ([]models.Expression, error) {
	rows, err := DB.Query("SELECT expression, status, result, time_start, time_finish FROM expressions")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, fmt.Errorf("Error querying database: %v", err)
	}
	defer rows.Close()

	var expressions []models.Expression

	for rows.Next() {
		var expression models.Expression
		if err := rows.Scan(
			&expression.Expression,
			&expression.Status,
			&expression.Result,
			&expression.CreatedAt,
			&expression.FinishedAt); err != nil {
			log.Println("Error scanning row:", err)
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		expressions = append(expressions, expression)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return expressions, nil
}

func FetchExpressionByID(id int) (*models.Expression, error) {
	row :=
		DB.QueryRow("SELECT expression, status, result, time_start, time_finish FROM expressions WHERE ID = ?", id)

	var expression models.Expression
	err := row.Scan(
		&expression.Expression,
		&expression.Status,
		&expression.Result,
		&expression.CreatedAt,
		&expression.FinishedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expression with ID %d not found", id)
		}
		log.Println("Error scanning row:", err)
		return nil, fmt.Errorf("Error scanning row: %v", err)
	}

	return &expression, nil
}
