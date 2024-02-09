package demonAndAgents

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"fmt"
	"log"
	"sync"
	"time"
)

var Servers = 4

func QueueHandler() {
	for {
		gotExpr, expression := queueMaster.ExpressionsQueue.Dequeue()
		if gotExpr {
			answerCh := make(chan bool)
			go ExpressionSeparator(expression, answerCh)
			<-answerCh
		}
	}
}

func ExpressionSeparator(expression models.Expression, answerCh chan bool) {
	fmt.Println(expression)

	// Здесь вы выполняете фактическую работу по обработке выражения
	// После завершения обработки отправьте сигнал в канал, чтобы сообщить, что работа завершена

	needCalculations := 10
	madeCalculations := 0

	ansCh := make(chan int, Servers)
	wg := &sync.WaitGroup{}

	operationTime, _ := cacheMaster.OperationCache.Get(cacheMaster.Operations["+"])

	for calc := 0; calc < needCalculations; calc++ {
		for id := 1; id <= Servers; id++ {
			if madeCalculations >= needCalculations {
				log.Println("finished calc")
				break
			}
			log.Println("added subcalc")
			madeCalculations++
			wg.Add(1)
			go CalculateSubExpression(id, expression.Expression, "+", operationTime, ansCh, wg)
		}
		wg.Wait()
	}

	time.Sleep(1 * time.Second)
	answerCh <- true
}

func CalculateSubExpression(id int, subExpression string, operation string, operationTime int, ansCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	server := models.Server{ID: id, Status: "Calculating SubExpression", Tasks: operation, LastPing: time.Now().Format("02-01-2006 15:04:05")}
	models.Servers.Mu.Lock()
	models.Servers.Servers[id] = &server
	models.Servers.Mu.Unlock()
	time.Sleep(time.Duration(operationTime) * time.Second)
	return
}
