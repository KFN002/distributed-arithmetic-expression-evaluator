package main

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/demonAgent"
	"distributed-arithmetic-expression-evaluator/backend/handlers"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
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

	log.Println("Перейти на сайт:", "http://localhost:8080/")

	data, err := databaseManager.ToCalculate()
	if err != nil {
		log.Println("Error fetching data from the database:", err)
		return
	}

	times, err := databaseManager.GetTimes()
	if err != nil {
		log.Println("Error fetching data from the database:", err)
		return
	}

	go cacheMaster.OperationCache.SetList(times)

	go queueMaster.ExpressionsQueue.EnqueueList(data) // загрузка в очередь выражений, которые мы не посчитали

	go demonAgent.QueueHandler() // начало работы обработчика данных

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
