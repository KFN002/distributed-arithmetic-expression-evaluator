package handlers

import (
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/utils"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	mu                sync.Mutex
	exampleExpression = models.Expression{
		ID:         0,
		Expression: "0",
		Status:     "example",
		CreatedAt:  time.Now().Format("02-01-2006 15:04:05"),
	}
)

// HandleAddExpression добавление выражения
func HandleAddExpression(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/create_expression.html")
	if err != nil {
		log.Println("Error parsing create_expression.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		input := r.FormValue("expression")

		input = strings.ReplaceAll(input, " ", "")

		validPattern := `^[0-9+\-*/()]+$`
		match, err := regexp.MatchString(validPattern, input)
		if err != nil {
			log.Println("Error in regular expression matching:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var status string

		if !match || !utils.CheckExpression(input) {
			status = "parsing error"
		} else {
			status = "processing"
		}

		double, err := databaseManager.CheckDuplicate(input)
		if err != nil {
			log.Println("Error fetching data", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var expression models.Expression
		expression.Expression = input
		expression.CreatedAt = time.Now().Format("02-01-2006 15:04:05")
		expression.Status = status

		if double {
			var fastExpression models.Expression
			fastExpression.ID, err = databaseManager.GetId(input)
			if err != nil {
				log.Println("Error fetching data", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			expression, err := databaseManager.FetchExpressionByID(fastExpression.ID)
			if err != nil {
				log.Println("Error fetching data", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			fastExpression.Status = expression.Status
			err = tmpl.Execute(w, fastExpression)
			if err != nil {
				log.Println("Error executing create_expression.html template:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}

		mu.Lock()
		defer mu.Unlock()

		result, err := databaseManager.DB.Exec("INSERT INTO expressions (expression, status, time_start) VALUES (?, ?, ?)",
			expression.Expression, expression.Status, expression.CreatedAt) // добавление бд в валидного выражения

		if err != nil {
			log.Println("Error inserting expression into database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		expressionID, err := result.LastInsertId()
		if err != nil {
			log.Println("Error getting last insert ID:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		expression.ID = int(expressionID)
		if expression.Status == "processing" { // добавление в очередь валидного выражения
			queueMaster.ExpressionsQueue.Enqueue(expression)
			log.Println("added to queue")
		}

		err = tmpl.Execute(w, expression)
		if err != nil {
			log.Println("Error executing create_expression.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error executing create_expression.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
