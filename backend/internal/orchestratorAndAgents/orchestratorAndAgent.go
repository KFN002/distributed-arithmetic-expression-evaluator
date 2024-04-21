package orchestratorAndAgents

import (
	"context"
	"errors"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/calculator"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/queueMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	pb "github.com/KFN002/distributed-arithmetic-expression-evaluator.git/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"sync"
	"time"
)

// QueueHandler Получение выражений из очереди
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
				if err := databaseManager.DB.UpdateExpressionAfterCalc(&expression); err != nil {
					log.Println("Error occurred when writing data:", err)
					queueMaster.ExpressionsQueue.Enqueue(expression)
				}

			case errCalc := <-errCh:
				log.Println(errCalc.Error())
				if errCalc.Error() == "division by zero or else" {
					log.Println("division 0")
					expression.ChangeData("calc error", 0)
					if err := databaseManager.DB.UpdateExpressionAfterCalc(&expression); err != nil {
						log.Println("Error occurred when writing data:", err)
						queueMaster.ExpressionsQueue.Enqueue(expression)
					}
				} else {
					log.Println("Error occurred:", errCalc)
					queueMaster.ExpressionsQueue.Enqueue(expression)
				}
			}
		}
	}
}

// Orchestrator Разделение выражения на подзадачи, подсчет выражения
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
					Agent(serversUsing, firstNum, secondNum, op, resCh, errSubCh, expression.UserID)
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

// Agent Подсчет мелкого выражения
func Agent(id int, firstNum float64, secondNum float64, operation string, subResCh chan float64, errCh chan error, userID int) {
	subExpression := fmt.Sprintf("%f %s %f", firstNum, operation, secondNum)
	log.Println(subExpression)

	go models.Servers.UpdateServers(id, subExpression, "Online, processing subExpression")

	operationID := cacheMaster.Operations[operation]
	operationTime, found := cacheMaster.OperationCache.Get(userID, operationID)
	if !found {
		errCh <- fmt.Errorf("operation time not found in cache for user ID %d and operation %s", userID, operation)
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
		// связь с gRPC сервером
		conn, err := grpc.Dial("localhost:8050", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("failed to dial server: %v", err)
		}
		defer conn.Close()

		grpcClient := pb.NewAgentServiceClient(conn)

		// запрос gRPC серверу
		result, err := grpcClient.Calculate(context.Background(), &pb.CalculationRequest{
			FirstNumber:  float32(firstNum),
			SecondNumber: float32(secondNum),
			Operation:    operation,
		})

		if result == nil {
			log.Fatal("grpcClient.Calculate returned nil result")
		}

		log.Println("Got gRPC response!")

		if err != nil {
			errCh <- errors.New("division by zero or else")
			go models.Servers.UpdateServers(id, subExpression, "Restarting, error occurred while processing")
			log.Println("calculating error - Agent calc")
		}

		go models.Servers.UpdateServers(id, "", "Online, finished processing")
		subResCh <- float64(result.Result)
	}
}
