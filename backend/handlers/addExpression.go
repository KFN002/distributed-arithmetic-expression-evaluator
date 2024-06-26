package handlers

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/queueMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/utils"
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

	userID := r.Context().Value("userID").(float64)
	login := r.Context().Value("login").(string)

	log.Println("User request:", userID, login)

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

		hasDuplicate, err := databaseManager.DB.CheckDuplicate(input, int(userID))
		if err != nil {
			log.Println("Error fetching data", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		expression := models.NewExpression(input, status, int(userID))

		if hasDuplicate {
			var newExpression models.Expression
			newExpression.ID, err = databaseManager.DB.GetId(input)
			if err != nil {
				log.Println("Error fetching data", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			originalExpression, err := databaseManager.DB.FetchExpressionByID(newExpression.ID, int(userID))
			if err != nil {
				log.Println("Error fetching data", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			newExpression.Status = originalExpression.Status
			err = tmpl.Execute(w, newExpression)
			if err != nil {
				log.Println("Error executing create_expression.html template:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}

		mu.Lock()
		defer mu.Unlock()

		result, err := databaseManager.DB.DB.Exec("INSERT INTO expressions (expression, status, time_start, user_id) VALUES (?, ?, ?, ?)",
			expression.Expression, expression.Status, expression.CreatedAt, userID) // добавление бд в валидного выражения

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
