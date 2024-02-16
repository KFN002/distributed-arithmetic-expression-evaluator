package orchestratorAndAgents

import (
	"distributed-arithmetic-expression-evaluator/backend/internal/cacheMaster"
	"distributed-arithmetic-expression-evaluator/backend/internal/calculator"
	"distributed-arithmetic-expression-evaluator/backend/internal/databaseManager"
	"distributed-arithmetic-expression-evaluator/backend/internal/queueMaster"
	"distributed-arithmetic-expression-evaluator/backend/pkg/models"
	"errors"
	"fmt"
	"log"
	"strconv"
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

	var answers []float64
	var freeServers = models.Servers.ServersQuantity
	var serversUsing = 0
	wg := &sync.WaitGroup{}

	for _, elem := range postfixExpression {
		if elem == "+" || elem == "-" || elem == "*" || elem == "/" {
			if freeServers > 0 {
				wg.Add(1)
				resCh := make(chan float64)
				errSubCh := make(chan error)
				firstNum := answers[len(answers)-2]
				secondNum := answers[len(answers)-1]
				answers = answers[:len(answers)-2]
				go func(firstNum, secondNum float64, op string) {
					defer wg.Done()
					Agent(serversUsing, firstNum, secondNum, op, resCh, errSubCh)
				}(firstNum, secondNum, elem)
				freeServers--
				serversUsing++
				select {
				case err := <-errSubCh:
					errCh <- err
					return
				case result := <-resCh:
					answers = append(answers, result)
				}
				freeServers++
				serversUsing--
			} else {
				wg.Wait()
			}
		} else {
			num, _ := strconv.Atoi(elem)
			answers = append(answers, float64(num))
		}
	}
	wg.Wait()
	answerCh <- answers[0]
}

func Agent(id int, firstNum float64, secondNum float64, operation string, subResCh chan float64, errCh chan error) {

	subExpression := fmt.Sprintf("%f %s %f", firstNum, operation, secondNum)
	log.Println(subExpression)

	models.Servers.UpdateServers(id, subExpression, "Online, processing subExpression")

	operationTime, err := cacheMaster.OperationCache.Get(cacheMaster.Operations[operation])
	if err != true {
		errCh <- errors.New("calculation error")
		models.Servers.UpdateServers(id, subExpression, "Restarting, error occurred while processing")
		log.Println("calculating error")
		return
	}

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

		result, err := calculator.Calculate(firstNum, secondNum, operation)
		if err != nil {
			errCh <- errors.New("calculation error")
			models.Servers.UpdateServers(id, subExpression, "Restarting, error occurred while processing")
			log.Println("calculating error")
			return
		}

		log.Println(result)

		models.Servers.UpdateServers(id, "", "Online, finished processing")
		subResCh <- result
	}
}
