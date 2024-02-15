package orchestratorAndAgents

import (
	"distributed-arithmetic-expression-evaluator/backend/internal/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/internal/calculator"
	"distributed-arithmetic-expression-evaluator/backend/internal/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/internal/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/pkg/models"
	"distributed-arithmetic-expression-evaluator/backend/pkg/utils"
	"errors"
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

	postfixExpression := calculator.InfixToPostfix(expression.Expression)
	log.Println(calculator.Solve(postfixExpression, cacheMaster.Operations))

	needCalculations := utils.CountOperators(expression.Expression)
	madeCalculations := 0

	var answers []float64

	ansCh := make(chan float64)
	wg := &sync.WaitGroup{}

	operationTime, _ := cacheMaster.OperationCache.Get(cacheMaster.Operations["+"])

	var calculated bool

	for calc := 0; calc < needCalculations; calc++ {
		for id := 1; id <= models.Servers.ServersQuantity; id++ {
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

	models.Servers.UpdateServers(id, subExpression, "Online, processing subExpression")

	subExpression = utils.RemoveRedundantParentheses(subExpression)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				models.Servers.SendHeartbeat(id)
			}
		}
	}()

	select {
	case <-time.After(time.Duration(operationTime) * time.Second):
		result, err := utils.CalculateSimpleTask(subExpression)
		if err != nil {
			errCh <- errors.New("calculation error")
			log.Println("calculating error")
			return
		}
		models.Servers.UpdateServers(id, "", "Online, finished processing")
		subResCh <- result
	}
}
