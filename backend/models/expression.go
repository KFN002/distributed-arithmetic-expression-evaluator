package models

import "time"

// Expression Арифметическое выражение.
type Expression struct {
	ID         int     `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     *int    `json:"result,omitempty"`
	CreatedAt  string  `json:"created_at"`
	FinishedAt *string `json:"finished_at,omitempty"`
}

func ChangeExpressionData(expression *Expression, status string, result int) {
	expression.Status = status
	expression.Result = &result
	finish := time.Now().Format("02-01-2006 15:04:05")
	expression.FinishedAt = &(finish)
}
