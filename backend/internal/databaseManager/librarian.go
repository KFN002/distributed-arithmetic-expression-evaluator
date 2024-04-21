package databaseManager

import (
	"database/sql"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var DB *DataBase

type DataBase struct {
	DB *sql.DB
	mu sync.Mutex
}

func init() {
	var err error
	DB, err = NewDataBase("backend/internal/database/database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func NewDataBase(dataSourceName string) (*DataBase, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DataBase{DB: db}, nil
}

func CloseDB() {
	if DB != nil {
		if err := DB.DB.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}
}

func (db *DataBase) GetTimes(userID int) ([]int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.DB.Query("SELECT time FROM operations WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var operationTimes []int
	for rows.Next() {
		var time int
		if err := rows.Scan(&time); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		operationTimes = append(operationTimes, time)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return operationTimes, nil
}

func (db *DataBase) CheckDuplicate(expression string, userID int) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM expressions WHERE expression = ? AND user_id = ?", expression, userID).Scan(&count)
	if err != nil {
		log.Println("Error querying database:", err)
		return false, fmt.Errorf("error querying database: %v", err)
	}
	return count > 0, nil
}

func (db *DataBase) GetId(expression string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var id int
	err := db.DB.QueryRow("SELECT id FROM expressions WHERE expression = ?", expression).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Println("Error querying database:", err)
		return 0, fmt.Errorf("error querying database: %v", err)
	}

	return id, nil
}

func (db *DataBase) GetExpressions(userID int) ([]models.Expression, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.DB.Query("SELECT * FROM expressions WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var expressions []models.Expression

	for rows.Next() {
		var expression models.Expression
		if err := rows.Scan(
			&expression.ID,
			&expression.Expression,
			&expression.Status,
			&expression.Result,
			&expression.CreatedAt,
			&expression.FinishedAt,
			&expression.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		expressions = append(expressions, expression)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return expressions, nil
}

func (db *DataBase) FetchExpressionByID(id, userID int) (*models.Expression, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.DB.QueryRow("SELECT * FROM expressions WHERE ID = ? AND user_id = ?", id, userID)

	var expression models.Expression
	err := row.Scan(
		&expression.ID,
		&expression.Expression,
		&expression.Status,
		&expression.Result,
		&expression.CreatedAt,
		&expression.FinishedAt,
		&expression.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expression with ID %d not found", id)
		}
		log.Println("Error scanning row:", err)
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	return &expression, nil
}

func (db *DataBase) ToCalculate() ([]models.Expression, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	rows, err := db.DB.Query("SELECT * FROM expressions WHERE status = ?", "processing")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var expressions []models.Expression

	for rows.Next() {
		var expression models.Expression
		if err := rows.Scan(
			&expression.ID,
			&expression.Expression,
			&expression.Status,
			&expression.Result,
			&expression.CreatedAt,
			&expression.FinishedAt,
			&expression.UserID); err != nil {
			log.Println("Error scanning row:", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		expressions = append(expressions, expression)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return expressions, nil
}

func (db *DataBase) UpdateExpressionAfterCalc(expression *models.Expression) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	stmt, err := db.DB.Prepare("UPDATE expressions SET status = ?, result = ?, time_finish = ? WHERE id = ?")
	if err != nil {
		log.Println("Error preparing update statement:", err)
		return fmt.Errorf("error preparing update statement: %v", err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("error closing file")
		}
	}(stmt)

	_, err = stmt.Exec(expression.Status, expression.Result, expression.FinishedAt, expression.ID)
	if err != nil {
		log.Println("Error updating expression:", err)
		return fmt.Errorf("error updating expression: %v", err)
	}

	return nil
}

func (db *DataBase) AddOperations(userID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	stmt, err := db.DB.Prepare("INSERT INTO operations (name, time, user_id) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing SQL statement: %v", err)
	}
	defer stmt.Close()

	operations := []models.Operations{
		{Name: "+", Time: 1, UserID: userID},
		{Name: "-", Time: 1, UserID: userID},
		{Name: "*", Time: 1, UserID: userID},
		{Name: "/", Time: 1, UserID: userID},
	}

	for _, op := range operations {
		_, err := stmt.Exec(op.Name, op.Time, op.UserID)
		if err != nil {
			return fmt.Errorf("error inserting operation into database: %v", err)
		}
	}

	return nil
}

func (db *DataBase) UpdateOperationTime(value int, name string, userID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.DB.Exec("UPDATE operations SET time=? WHERE name=? AND user_id=?", value, name, userID)
	if err != nil {
		return fmt.Errorf("error updating operation time: %v", err)
	}
	return nil
}

func (db *DataBase) GetUserIDs() ([]int, error) {
	var userIDs []int

	rows, err := db.DB.Query("SELECT DISTINCT id FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying user IDs from the database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning user ID row: %v", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user ID rows: %v", err)
	}

	return userIDs, nil
}
