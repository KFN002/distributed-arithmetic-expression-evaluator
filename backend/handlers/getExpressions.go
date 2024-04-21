package handlers

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/utils"
	"html/template"
	"log"
	"net/http"
)

// HandleExpressions возвращает страницу с данными выражений
func HandleExpressions(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(float64)

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

	expressions, err := databaseManager.DB.GetExpressions(int(userID)) // получаем выражения
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
