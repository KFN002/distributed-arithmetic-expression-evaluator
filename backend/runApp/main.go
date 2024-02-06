package main

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
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
	switch r.Method {
	case http.MethodPost:
		handleAddExpression(w, r)
	case http.MethodGet:
		getExpressionsHTML(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getExpressionsHTML(w http.ResponseWriter, r *http.Request) {
	// Assuming you have an expressions.html template
	tmpl, err := template.ParseFiles("static/assets/view_expressions.html")
	if err != nil {
		log.Println("Error parsing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	//place expressions into template

	err = tmpl.Execute(w, expressions)
	if err != nil {
		log.Println("Error executing expressions.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleChangeCalcTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// get new data
		// update the database
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
	if r.Method != http.MethodGet {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var expression models.Expression
		if err := json.Unmarshal(body, &expression); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		expression.ID = 1
		expression.CreatedAt = time.Now()
		expression.Status = "processing"

		mu.Lock()
		expressions = append(expressions, expression)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expression)

		// get expression
		// add expression to a database (sync)
		// add expression to a redis queue
		// check if expression has occurred previously

		http.Redirect(w, r, "/expressions", 200)
	}

	// Assuming you have an add-expression.html template
	tmpl, err := template.ParseFiles("static/assets/create_expression.html")
	if err != nil {
		log.Println("Error parsing create_expression.html template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	err = tmpl.Execute(w, nil)
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
