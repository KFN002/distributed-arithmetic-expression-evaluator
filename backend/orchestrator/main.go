package main

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
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
		input1 := r.FormValue("time1")
		input2 := r.FormValue("time2")
		input3 := r.FormValue("time3")
		input4 := r.FormValue("time4")

		fmt.Println(input1, input2, input3, input4)

		return
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
		var expression models.Expression
		input := r.FormValue("expression")

		fmt.Println(input)

		expression.ID = 1 // Assuming the ID is 1 for simplicity
		expression.Expression = input
		expression.CreatedAt = time.Now().Format("02-01-2006 15:04:05")
		expression.Status = "processing"

		// add expression to a database (sync)
		// add expression to a redis queue
		// check if expression has occurred previously

		mu.Lock()
		expressions = append(expressions, expression)
		mu.Unlock()

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
