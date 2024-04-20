package tests

import (
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/cacheMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/orchestratorAndAgents"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"log"
	"sync"
	"testing"
)

func TestOrchestrator(t *testing.T) {
	expression := models.Expression{
		Expression: "2 + 3 * 4",
		UserID:     1,
	}

	answerCh := make(chan float64)
	errCh := make(chan error)

	cacheMaster.OperationCache.Set(expression.UserID, cacheMaster.Operations["+"], 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		orchestratorAndAgents.Orchestrator(expression, answerCh, errCh)
	}()

	select {
	case result := <-answerCh:
		fmt.Println("Result:", result)
	case err := <-errCh:
		log.Println("Error:", err)
	}

	wg.Wait()
}
