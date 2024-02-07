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

type Server struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}
