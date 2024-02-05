package main

import "time"

// Expression представляет арифметическое выражение.
type Expression struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	Status    string    `json:"status"`
	Result    string    `json:"result,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
