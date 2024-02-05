package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	expressions []Expression
	mu          sync.Mutex
)

// handleExpressions обрабатывает запросы по управлению арифметическими выражениями.
func handleExpressions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addExpression(w, r)
	case http.MethodGet:
		getExpressions(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// addExpression добавляет новое арифметическое выражение.
func addExpression(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var expression Expression
	if err := json.Unmarshal(body, &expression); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	expression.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	expression.CreatedAt = time.Now()
	expression.Status = "processing"

	mu.Lock()
	expressions = append(expressions, expression)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expression)
}

// getExpressions возвращает список арифметических выражений.
func getExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(expressions)
}
