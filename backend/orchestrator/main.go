package main

import (
	"distributed-arithmetic-expression-evaluator/backend/dataManager"
	_ "distributed-arithmetic-expression-evaluator/backend/dataManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	expressions []models.Expression
	mu          sync.Mutex
)

func handleExpressions(w http.ResponseWriter, r *http.Request) {
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

	//place expressions into template
	if len(expressions) == 0 {
		expressions = append(expressions,
			models.Expression{ID: 0,
				Expression: "0",
				Status:     "example",
				Result:     0,
				CreatedAt:  time.Now().Format("02-01-2006 15:04:05"),
				FinishedAt: time.Now().Format("02-01-2006 15:04:05")})
	}

	err = tmpl.Execute(w, expressions)
	if err != nil {
		log.Println("Error executing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// get new data
		// update the database

		timeAdd := r.FormValue("time1")
		timeSub := r.FormValue("time2")
		timeMult := r.FormValue("time3")
		timeDiv := r.FormValue("time4")

		time1, err := strconv.Atoi(timeAdd)
		if err != nil {
			http.Error(w, "Invalid input for time1", http.StatusBadRequest)
			return
		}
		time2, err := strconv.Atoi(timeSub)
		if err != nil {
			http.Error(w, "Invalid input for time2", http.StatusBadRequest)
			return
		}
		time3, err := strconv.Atoi(timeMult)
		if err != nil {
			http.Error(w, "Invalid input for time3", http.StatusBadRequest)
			return
		}
		time4, err := strconv.Atoi(timeDiv)
		if err != nil {
			http.Error(w, "Invalid input for time4", http.StatusBadRequest)
			return
		}

		if time1 <= 0 || time2 <= 0 || time3 <= 0 || time4 <= 0 {
			http.Error(w, "Input values must be positive and not zero", http.StatusBadRequest)
			return
		}

		fmt.Println(timeAdd, timeSub, timeMult, timeDiv)

		http.Redirect(w, r, "/change-calc-time", 200)
	}

	// Assuming you have a change-calc-time.html template
	tmpl, err := template.ParseFiles("static/assets/edit_time.html")
	if err != nil {
		log.Println("Error parsing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error executing edit_time.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleAddExpression(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/assets/create_expression.html")
	if err != nil {
		log.Println("Error parsing create_expression.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		input := r.FormValue("expression")

		// Regular expression pattern to match only integers and allowed operators
		validPattern := "^[0-9+\\-*/()\\s]+$"
		match, err := regexp.MatchString(validPattern, input)
		if err != nil {
			log.Println("Error in regular expression matching:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if !match {
			http.Error(w, "Expression contains invalid characters", http.StatusBadRequest)
			return
		}

		input = strings.ReplaceAll(input, " ", "") // Remove spaces from input

		// Proceed with handling the expression
		var expression models.Expression
		expression.Expression = input
		expression.CreatedAt = time.Now().Format("02-01-2006 15:04:05")
		expression.Status = "processing"

		// Insert expression into the database
		result, err := dataManager.
			DB.Exec("INSERT INTO expressions (expression, status, time_start) VALUES (?, ?, ?)",
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

func handleCurrentServers(w http.ResponseWriter, r *http.Request) {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handleExpressions)
	r.HandleFunc("/expressions", handleExpressions)
	r.HandleFunc("/change-calc-time", handleChangeCalcTime) // Implement this
	r.HandleFunc("/add-expression", handleAddExpression)    // Implement this
	r.HandleFunc("/current-servers", handleCurrentServers)  // Implement this

	// Serve static files
	fileServer := http.FileServer(http.Dir("static/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
