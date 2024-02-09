package models

// Арифметическое выражение.
type Expression struct {
	ID         int     `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     *int    `json:"result,omitempty"`
	CreatedAt  string  `json:"created_at"`
	FinishedAt *string `json:"finished_at,omitempty"`
}

// Данные сервера
type Server struct {
	ID       int    `json:"id"`
	Tasks    string `json:"task"`
	LastPing int    `json:"ping"`
}

// Данные операций - время
type OperationTimes struct {
	Time1 int
	Time2 int
	Time3 int
	Time4 int
}

// общие данные о серверах
type ServersData struct {
	Servers map[int]*Server
}

// стак и его логика
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
