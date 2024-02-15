package orchestratorAndAgents

import (
	"distributed-arithmetic-expression-evaluator/backend/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/utils"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

func QueueHandler() {
	for {
		expression, gotExpr := queueMaster.ExpressionsQueue.Dequeue()
		if gotExpr {
			answerCh := make(chan float64)
			errCh := make(chan error)

			go Orchestrator(expression, answerCh, errCh)

			select {
			case ans := <-answerCh:

				expression.ChangeData("finished", ans)
				if err := databaseManager.UpdateExpressionAfterCalc(&expression); err != nil {
					log.Println("Error occurred when writing data:", err)
					queueMaster.ExpressionsQueue.Enqueue(expression)
				}
				log.Println(ans)

			case err := <-errCh:
				log.Println("Error occurred:", err)
				queueMaster.ExpressionsQueue.Enqueue(expression)
			}
		}
	}
}

func Orchestrator(expression models.Expression, answerCh chan float64, errCh chan error) {
	defer close(answerCh)
	fmt.Println(expression)

	// Здесь вы выполняете фактическую работу по обработке выражения
	// После завершения обработки отправьте сигнал в канал, чтобы сообщить, что работа завершена

	needCalculations := utils.CountOperators(expression.Expression)
	madeCalculations := 0

	var answers []float64

	ansCh := make(chan float64)
	wg := &sync.WaitGroup{}

	operationTime, _ := cacheMaster.OperationCache.Get(cacheMaster.Operations["+"])

	var calculated bool

	for calc := 0; calc < needCalculations; calc++ {
		for id := 1; id <= models.ServersQuantity; id++ {
			if madeCalculations >= needCalculations {
				log.Println("finished calc")
				calculated = true
				break
			}
			log.Println("added subcalc")

			madeCalculations++

			wg.Add(1)

			go Agent(id, expression.Expression, operationTime, ansCh, errCh, wg)
		}

		if calculated {
			break
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

func Agent(id int, subExpression string, operationTime int, subResCh chan float64, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println(operationTime)

	models.UpdateServers(id, subExpression, "Online, processing subExpression")

	subExpression = utils.RemoveRedundantParentheses(subExpression)

	select {
	case <-time.After(time.Duration(operationTime) * time.Second):
		result, err := utils.CalculateSimpleTask(subExpression)
		if err != nil {
			errCh <- errors.New("calculation error")
			log.Println("calculating error")
			return
		}
		models.UpdateServers(id, "", "Online, finished processing")
		subResCh <- result
	}
}
