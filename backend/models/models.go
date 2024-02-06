package models

import "time"

// Expression представляет арифметическое выражение.
type Expression struct {
	ID         int       `json:"id"`
	Status     string    `json:"status"`
	Result     int       `json:"result,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

// Result представляет результат вычисления арифметического выражения.
type Result struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

// Operation представляет операцию с временем её выполнения.
type Operation struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}
