package models

// Expression представляет арифметическое выражение.
type Expression struct {
	ID         int    `json:"id"`
	Expression string `json:"expression"`
	Status     string `json:"status"`
	Result     int    `json:"result,omitempty"`
	CreatedAt  string `json:"created_at"`
	FinishedAt string `json:"finished_at,omitempty"`
}

// Result представляет результат вычисления арифметического выражения.
type Result struct {
	ID     string `json:"id"`
	Result string `json:"result"`
}

// Operation представляет операцию с временем её выполнения.
type Operation struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}
