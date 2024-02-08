package main

import (
	"distributed-arithmetic-expression-evaluator/backend/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HandleExpressions)
	r.HandleFunc("/expressions", handlers.HandleExpressions)
	r.HandleFunc("/change-calc-time", handlers.HandleChangeCalcTime)
	r.HandleFunc("/add-expression", handlers.HandleAddExpression)
	r.HandleFunc("/current-servers", handlers.HandleCurrentServers)
	r.HandleFunc("/expression-by-id", handlers.HandleGetExpressionByID)

	fileServer := http.FileServer(http.Dir("static/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	fmt.Println("Перейти на сайт:", "http://localhost:8080/")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
