package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"testing"
	"time"
)

func TestChangeData(t *testing.T) {
	expr := models.Expression{
		ID:         1,
		Expression: "2+3",
		Status:     "pending",
		Result:     nil,
		CreatedAt:  "01-01-2022 10:00:00",
		FinishedAt: nil,
		UserID:     1,
	}

	expectedStatus := "completed"
	expectedResult := 5.0
	expr.ChangeData(expectedStatus, expectedResult)

	if expr.Status != expectedStatus {
		t.Errorf("ChangeData() failed to update status. Got: %s, Expected: %s", expr.Status, expectedStatus)
	}

	if *expr.Result != expectedResult {
		t.Errorf("ChangeData() failed to update result. Got: %f, Expected: %f", *expr.Result, expectedResult)
	}

	currentTime := time.Now().Format("02-01-2006 15:04:05")
	if *expr.FinishedAt != currentTime {
		t.Errorf("ChangeData() failed to update finished time. Got: %s, Expected: %s", *expr.FinishedAt, currentTime)
	}
}

func TestNewExpression(t *testing.T) {
	expression := "2*3"
	status := "pending"
	userID := 1
	startTime := time.Now().Format("02-01-2006 15:04:05")
	expr := models.NewExpression(expression, status, userID)

	if expr.Expression != expression {
		t.Errorf("NewExpression() failed to set expression. Got: %s, Expected: %s", expr.Expression, expression)
	}

	if expr.Status != status {
		t.Errorf("NewExpression() failed to set status. Got: %s, Expected: %s", expr.Status, status)
	}

	if expr.UserID != userID {
		t.Errorf("NewExpression() failed to set userID. Got: %d, Expected: %d", expr.UserID, userID)
	}

	if expr.CreatedAt != startTime {
		t.Errorf("NewExpression() failed to set created time. Got: %s, Expected: %s", expr.CreatedAt, startTime)
	}
}
