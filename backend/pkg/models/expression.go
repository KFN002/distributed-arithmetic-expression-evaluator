package models

import "time"

// Expression выражение
type Expression struct {
	ID         int
	Expression string
	Status     string
	Result     *float64
	CreatedAt  string
	FinishedAt *string
	UserID     int
}

// ChangeData Изменение данных выражения
func (e *Expression) ChangeData(status string, result float64) {
	e.Status = status
	e.Result = &result
	finish := time.Now().Format("02-01-2006 15:04:05")
	e.FinishedAt = &finish
}

// NewExpression создание нового экземпляра класса выражение
func NewExpression(expression string, status string, userID int) Expression {
	var e Expression
	e.Status = status
	e.Expression = expression
	start := time.Now().Format("02-01-2006 15:04:05")
	e.CreatedAt = start
	e.UserID = userID
	return e
}
