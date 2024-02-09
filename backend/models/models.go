package models

import (
	"sync"
	"time"
)

var Servers = ServersData{Servers: map[int]*Server{
	1: {ID: 1, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	2: {ID: 2, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	3: {ID: 3, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
	4: {ID: 4, Status: "Stand By", Tasks: "", LastPing: time.Now().Format("02-01-2006 15:04:05")},
}}

// Expression Арифметическое выражение.
type Expression struct {
	ID         int     `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     *int    `json:"result,omitempty"`
	CreatedAt  string  `json:"created_at"`
	FinishedAt *string `json:"finished_at,omitempty"`
}

// Server Данные сервера
type Server struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Tasks    string `json:"task"`
	LastPing string `json:"ping"`
}

// ServersData данные о серверах
type ServersData struct {
	Mu      sync.Mutex
	Servers map[int]*Server
}

// OperationTimes Данные операций - время
type OperationTimes struct {
	Time1 int
	Time2 int
	Time3 int
	Time4 int
}

// Stack стак и его логика
type Stack []string

func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

func (st *Stack) Push(str string) {
	*st = append(*st, str)
}

func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1
		*st = (*st)[:index]
		return true
	}
}

func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1
		element := (*st)[index]
		return element
	}
}
