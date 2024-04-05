package main

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/handlers"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/orchestratorAndAgents"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/queueMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", JWTMiddleware(handlers.HandleExpressions))
	r.HandleFunc("/expressions", JWTMiddleware(handlers.HandleExpressions))
	r.HandleFunc("/change-calc-time", JWTMiddleware(handlers.HandleChangeCalcTime))
	r.HandleFunc("/add-expression", JWTMiddleware(handlers.HandleAddExpression))
	r.HandleFunc("/current-servers", JWTMiddleware(handlers.HandleCurrentServers))
	r.HandleFunc("/expression-by-id", JWTMiddleware(handlers.HandleGetExpressionByID))
	r.HandleFunc("/scheme", JWTMiddleware(handlers.HandleGetScheme))

	r.HandleFunc("/login", handlers.HandleLogin)
	r.HandleFunc("/signup", handlers.HandleRegister)

	fileServer := http.FileServer(http.Dir("static/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	log.Println("Перейти в интерфейс:", "http://localhost:8080/")

	// подгрузка заданий для калькуляции
	data, err := databaseManager.ToCalculate()
	if err != nil {
		log.Println("Error fetching data from the database:", err)
		return
	}

	// подгрузка времени выполнения операции
	times, err := databaseManager.GetTimes()
	if err != nil {
		log.Println("Error fetching data from the database:", err)
		return
	}

	// подключение "серверов"
	go models.Servers.InitServers()
	go models.Servers.RunServers()

	// загрузка в кэш данных об операциях, чтобы не делать запрос в бд каждый раз
	go cacheMaster.OperationCache.SetList(times)

	go queueMaster.ExpressionsQueue.EnqueueList(data) // загрузка в очередь выражений, которые мы не посчитали

	go orchestratorAndAgents.QueueHandler() // начало работы обработчика данных - постоянно читает данные из очереди

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
