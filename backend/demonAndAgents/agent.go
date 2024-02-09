package demonAndAgents

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/utils"
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
			answerCh := make(chan int)

			go ExpressionSeparator(expression, answerCh)
			ans := <-answerCh

			models.ChangeExpressionData(&expression, "finished", ans)
			if err := databaseManager.UpdateExpressionAfterCalc(&expression); err != nil {
				log.Println("Error occurred when writing data:", err)
				models.ChangeExpressionData(&expression, "failed", ans)
				queueMaster.ExpressionsQueue.Enqueue(expression)
			}

			log.Println(ans)
		}
	}
}

func ExpressionSeparator(expression models.Expression, answerCh chan int) {
	defer close(answerCh)
	fmt.Println(expression)

	// Здесь вы выполняете фактическую работу по обработке выражения
	// После завершения обработки отправьте сигнал в канал, чтобы сообщить, что работа завершена

	needCalculations := 10
	madeCalculations := 0

	var answers []int

	ansCh := make(chan int)
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
	}

	go func() {
		wg.Wait()
		close(ansCh)
	}()

	for res := range ansCh {
		answers = append(answers, res)
	}

	answerCh <- utils.SumList(answers)
}

func CalculateSubExpression(id int, subExpression string, operation string, operationTime int, subResCh chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	updateServerStatus := func(status string) {
		models.UpdateServers(id, subExpression, status)
	}

	updateServerStatus("Online, processing subExpression")

	timer := time.After(time.Duration(operationTime) * time.Second)

	select {
	case <-timer:
		updateServerStatus("Restarting, calculation failed")
		return
	default:
		result := 5
		updateServerStatus("Online, finished processing")
		subResCh <- result
	}
}
