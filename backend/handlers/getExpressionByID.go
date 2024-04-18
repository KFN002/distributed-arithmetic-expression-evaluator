package handlers

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// HandleGetExpressionByID получение выражения по id
func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(float64)

	tmpl, err := template.ParseFiles("static/assets/expression_by_id.html")
	if err != nil {
		log.Println("Error parsing expression_by_id.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		expression, err := databaseManager.DB.FetchExpressionByID(id, int(userID))
		if err != nil {
			http.Error(w, "Failed to fetch expression", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, expression)
		if err != nil {
			log.Println("Error executing expression_by_id.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		err = tmpl.Execute(w, exampleExpression)
		if err != nil {
			log.Println("Error executing expression_by_id.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
