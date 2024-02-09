package agent

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"fmt"
	"sync"
	"time"
)

var Servers = 3

func QueueHandler() {
	for {
		gotExpr, expression := queueMaster.ExpressionsQueue.Dequeue()
		if gotExpr {
			answerCh := make(chan bool)
			ExpressionSeparator(expression, answerCh)
			<-answerCh
		}
	}
}

func ExpressionSeparator(expression models.Expression, answerCh chan bool) {
	fmt.Println(expression)

	// Здесь вы выполняете фактическую работу по обработке выражения
	// После завершения обработки отправьте сигнал в канал, чтобы сообщить, что работа завершена
	ansCh := make(chan int, Servers)
	wg := &sync.WaitGroup{}

	operationTime, _ := cacheMaster.OperationCache.Get(cacheMaster.Operations["+"])

	for i := 0; i < Servers; i++ {
		wg.Add(1)
		go CalculateSubExpression(i, expression.Expression, "+", operationTime, ansCh, wg)
	}

	wg.Wait()

	time.Sleep(1 * time.Second)
	answerCh <- true
}

func CalculateSubExpression(id int, subExpression string, operation string, operationTime int, ansCh chan<- int, wg *sync.WaitGroup) {
	wg.Done()
}
