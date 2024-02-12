package handlers

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	_ "distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
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

// HandleChangeCalcTime страница изменения времени выполнения
func HandleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		timeValues := map[string]int{}
		timeFields := []string{"1", "2", "3", "4"}
		for _, field := range timeFields {
			timeStr := r.FormValue(field)
			timeValue, err := strconv.Atoi(timeStr)
			if err != nil || timeValue <= 0 {
				http.Error(w, fmt.Sprintf("Invalid input for %s", field), http.StatusBadRequest)
				return
			}
			timeValues[field] = timeValue
		}

		for id, value := range timeValues {
			_, err := databaseManager.DB.Exec("UPDATE operations SET time=? WHERE id=?", value, id)
			if err != nil {
				log.Println("Error updating database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/change-calc-time", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("static/assets/edit_time.html")
	if err != nil {
		log.Println("Error parsing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationTimes, err := databaseManager.GetTimes() // получение данных время операций
	if err != nil {
		log.Println("Error fetching times", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationsData := models.OperationTimes{
		Time1: operationTimes[0],
		Time2: operationTimes[1],
		Time3: operationTimes[2],
		Time4: operationTimes[3],
	}

	go cacheMaster.OperationCache.SetList(operationTimes)

	err = tmpl.Execute(w, operationsData)
	if err != nil {
		log.Println("Error executing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

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

		validPattern := `^[0-9+\-*/()\s]*[^a-zA-Z!@#$%^&*_=<>?|\\.,;:~"']{2}[0-9+\-*/()\s]*$`
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

		input = strings.ReplaceAll(input, " ", "")

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

// HandleCurrentServers получение данных о серверах
func HandleCurrentServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("static/assets/server_data.html")
	if err != nil {
		log.Println("Error parsing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	models.Servers.Mu.Lock()

	var serverList []*models.Server

	for _, server := range models.Servers.Servers {
		serverList = append(serverList, server)
	}

	// сортировка серверов по id
	sort.Slice(serverList, func(i, j int) bool {
		return serverList[i].ID < serverList[j].ID
	})

	models.Servers.Mu.Unlock()

	err = tmpl.Execute(w, serverList)
	if err != nil {
		log.Println("Error executing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandleGetExpressionByID получение выражения по id
func HandleGetExpressionByID(w http.ResponseWriter, r *http.Request) {

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

		expression, err := databaseManager.FetchExpressionByID(id)
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

func HandleGetScheme(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/scheme.html")
	if err != nil {
		log.Println("Error parsing expression_by_id.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if r.Method == http.MethodGet {
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Error executing expression_by_id.html template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
