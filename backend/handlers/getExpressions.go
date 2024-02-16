package handlers

import (
	"distributed-arithmetic-expression-evaluator/backend/internal/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/pkg/utils"
	"html/template"
	"log"
	"net/http"
)

// HandleExpressions возвращает страницу с данными выражений
func HandleExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("static/assets/view_expressions.html")
	if err != nil {
		log.Println("Error parsing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	expressions, err := databaseManager.GetExpressions() // получаем выражения
	if err != nil {
		log.Println("Error getting expressions:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(expressions) == 0 {
		expressions = append(expressions, exampleExpression) // если выражений нет, то передадим пример
	}

	err = tmpl.Execute(w, utils.FlipList(expressions))
	if err != nil {
		log.Println("Error executing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
