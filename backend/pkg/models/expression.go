package models

import "time"

// Expression выражение
type Expression struct {
	ID         int      `json:"id"`
	Expression string   `json:"expression"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result,omitempty"`
	CreatedAt  string   `json:"created_at"`
	FinishedAt *string  `json:"finished_at,omitempty"`
}

// ChangeData Изменение данных выражения
func (e *Expression) ChangeData(status string, result float64) {
	e.Status = status
	e.Result = &result
	finish := time.Now().Format("02-01-2006 15:04:05")
	e.FinishedAt = &finish
}

func NewExpression(expression string, status string) Expression {
	var e Expression
	e.Status = status
	e.Expression = expression
	start := time.Now().Format("02-01-2006 15:04:05")
	e.CreatedAt = start
	return e
}
