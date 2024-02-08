package handlers

import (
	"distributed-arithmetic-expression-evaluator/backend/dataManager"
	_ "distributed-arithmetic-expression-evaluator/backend/dataManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	expressions       []models.Expression
	mu                sync.Mutex
	exampleExpression = models.Expression{
		ID:         0,
		Expression: "0",
		Status:     "example",
		CreatedAt:  time.Now().Format("02-01-2006 15:04:05"),
	}
)

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

	expressions, err := dataManager.GetExpressions()
	if err != nil {
		log.Println("Error getting expressions:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If no expressions are found, create an example expression
	if len(expressions) == 0 {
		expressions = append(expressions, exampleExpression)
	}

	// Execute the template with expressions data
	err = tmpl.Execute(w, utils.FlipList(expressions))
	if err != nil {
		log.Println("Error executing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

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

		// Update database with time values
		for id, value := range timeValues {
			_, err := dataManager.DB.Exec("UPDATE operations SET time=? WHERE id=?", value, id)
			if err != nil {
				log.Println("Error updating database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/change-calc-time", http.StatusSeeOther)
	}

	// Assuming you have a change-calc-time.html template
	tmpl, err := template.ParseFiles("static/assets/edit_time.html")
	if err != nil {
		log.Println("Error parsing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	operationTimes, err := dataManager.GetTimes()
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

	err = tmpl.Execute(w, operationsData)
	if err != nil {
		log.Println("Error executing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func HandleAddExpression(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/create_expression.html")
	if err != nil {
		log.Println("Error parsing create_expression.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		input := r.FormValue("expression")

		// Regular expression pattern to match only integers and allowed operators
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

		input = strings.ReplaceAll(input, " ", "") // Remove spaces from input

		double, err := dataManager.CheckDuplicate(input)
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
			fastExpression.ID, err = dataManager.GetId(input)
			if err != nil {
				log.Println("Error fetching data", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			fastExpression.Status = status
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

		if expression.Status == "processing" {
			fmt.Println("added to queue")
		}

		result, err := dataManager.DB.Exec("INSERT INTO expressions (expression, status, time_start) VALUES (?, ?, ?)",
			expression.Expression, expression.Status, expression.CreatedAt)

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

		// Pass the expression ID to the template
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

func HandleCurrentServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// Assuming you have a server_data.html template
	tmpl, err := template.ParseFiles("static/assets/server_data.html")
	if err != nil {
		log.Println("Error parsing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//get servers data from databases
	//place it in the .Execute method
	//place servers into template
	serversData := []struct {
		ID            int
		Name          string
		Status        string
		LastOperation string
	}{
		{1, "Server 1", "Online", "Backup"},
		{2, "Server 2", "Offline", "Restart"},
		{3, "Server 3", "Online", "Update"},
	}

	err = tmpl.Execute(w, serversData)
	if err != nil {
		log.Println("Error executing server_data.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

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

		expression, err := dataManager.FetchExpressionByID(id)
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
