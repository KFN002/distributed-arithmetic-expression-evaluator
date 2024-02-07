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

type OperationTimes struct {
	Time1 int
	Time2 int
	Time3 int
	Time4 int
}
