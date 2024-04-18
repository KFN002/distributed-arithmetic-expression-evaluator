package main

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/handlers"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/orchestratorAndAgents"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/queueMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/middleware"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", middleware.JWTMiddleware(handlers.HandleExpressions))
	r.HandleFunc("/expressions", middleware.JWTMiddleware(handlers.HandleExpressions))
	r.HandleFunc("/change-calc-time", middleware.JWTMiddleware(handlers.HandleChangeCalcTime))
	r.HandleFunc("/add-expression", middleware.JWTMiddleware(handlers.HandleAddExpression))
	r.HandleFunc("/current-servers", middleware.JWTMiddleware(handlers.HandleCurrentServers))
	r.HandleFunc("/expression-by-id", middleware.JWTMiddleware(handlers.HandleGetExpressionByID))
	r.HandleFunc("/scheme", middleware.JWTMiddleware(handlers.HandleGetScheme))
	r.HandleFunc("/logout", middleware.JWTMiddleware(handlers.HandleLogout))

	r.HandleFunc("/login", handlers.HandleLogin)
	r.HandleFunc("/signup", handlers.HandleRegister)

	fileServer := http.FileServer(http.Dir("static/assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	log.Println("Перейти в интерфейс:", "http://localhost:8080/")

	// подгрузка заданий для калькуляции
	data, err := databaseManager.DB.ToCalculate()
	if err != nil {
		log.Println("Expression error")
		log.Println("Error fetching data from the database:", err)
		return
	}

	err = cacheMaster.LoadOperationTimesIntoCache()
	if err != nil {
		log.Println("Error while caching data and updating user info:", err)
		return
	}

	log.Println(cacheMaster.OperationCache)

	// подключение "серверов"
	go models.Servers.InitServers()
	go models.Servers.RunServers()

	go queueMaster.ExpressionsQueue.EnqueueList(data) // загрузка в очередь выражений, которые мы не посчитали

	go orchestratorAndAgents.QueueHandler() // начало работы обработчика данных - постоянно читает данные из очереди

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
