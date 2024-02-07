package main

import (
	"distributed-arithmetic-expression-evaluator/backend/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HandleExpressions)
	r.HandleFunc("/expressions", handlers.HandleExpressions)
	r.HandleFunc("/change-calc-time", handlers.HandleChangeCalcTime) // Implement this
	r.HandleFunc("/add-expression", handlers.HandleAddExpression)    // Implement this
	r.HandleFunc("/current-servers", handlers.HandleCurrentServers)  // Implement this

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
