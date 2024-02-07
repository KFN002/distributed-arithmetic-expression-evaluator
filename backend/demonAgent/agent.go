package main

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	workerCount   int
	expressionsCh = make(chan models.Expression)
	wg            sync.WaitGroup
)

func main() {
	// Получаем количество рабочих из переменной окружения
	workerCountStr := os.Getenv("WORKER_COUNT")
	if workerCountStr == "" {
		log.Fatal("WORKER_COUNT is not set")
	}
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil {
		log.Fatalf("Failed to parse WORKER_COUNT: %v", err)
	}

	// Запускаем рабочие горутины
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i)
	}

	// Запускаем горутину для получения выражений
	go fetchExpressions()

	// Ожидаем завершения всех рабочих горутин
	wg.Wait()
}

func worker(id int) {
	defer wg.Done()
	// Выполняем вычисление
	fmt.Printf("ok")
}

func fetchExpressions() {
	for {
		time.Sleep(1 * time.Second) // Периодически получаем выражения
		resp, err := http.Get("http://localhost:8080/expressions")
		if err != nil {
			log.Printf("Failed to fetch expressions: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to fetch expressions: %s\n", resp.Status)
			continue
		}

		var expressions []models.Expression
		if err := json.NewDecoder(resp.Body).Decode(&expressions); err != nil {
			log.Printf("Failed to decode expressions: %v\n", err)
			continue
		}

		// Отправляем полученные выражения на канал expressions
		for _, expr := range expressions {
			expressionsCh <- expr
		}
	}
}

func calculate(expr string) string {
	// Разбиваем строку на подвыражения
	subExprs := strings.Split(expr, "+")

	// Создаем канал для результатов подвыражений
	resultCh := make(chan interface{}, len(subExprs))

	// Запускаем горутины для вычисления подвыражений
	for _, subExpr := range subExprs {
		go func(s string) {
			result, err := eval(s)
			if err != nil {
				log.Printf("Error evaluating expression %s: %v", s, err)
				resultCh <- 0 // Если возникла ошибка, отправляем 0
				return
			}
			resultCh <- result
		}(strings.TrimSpace(subExpr))
	}

	// Суммируем результаты подвыражений
	total := 0
	for i := 0; i < len(subExprs); i++ {
		result := <-resultCh
		switch v := result.(type) {
		case int:
			total += v
		case float64:
			total += int(v)
		}
	}

	return strconv.Itoa(total)
}

func eval(expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}

	switch v := result.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("unsupported result type: %T", v)
	}
}
