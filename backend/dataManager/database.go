package dataManager

import (
	"database/sql"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

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
	query := "SELECT COUNT(*) FROM expressions WHERE expression = ?"
	var count int
	err := DB.QueryRow(query, expression).Scan(&count)
	if err != nil {
		log.Println("Error querying database:", err)
		return false, fmt.Errorf("Error querying database: %v", err)
	}
	return count > 0, nil
}

func GetId(expression string) (int, error) {
	// Prepare the SQL query to retrieve the ID associated with the expression
	query := "SELECT id FROM expressions WHERE expression = ?"

	var id int
	err := DB.QueryRow(query, expression).Scan(&id)
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
	query := "SELECT * FROM expressions"

	// Execute the query to fetch all expressions
	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, fmt.Errorf("Error querying database: %v", err)
	}
	defer rows.Close()

	// Initialize a slice to hold the expressions
	var expressions []models.Expression

	// Iterate over the result set and populate the expressions slice
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

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	// Return the fetched expressions
	return expressions, nil
}
