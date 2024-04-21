package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/queueMaster"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"testing"
)

func TestEnqueueDequeue(t *testing.T) {
	q := queueMaster.ExpressionQueue()

	q.Enqueue(models.Expression{Expression: "1+2"})
	q.Enqueue(models.Expression{Expression: "3*4"})
	q.Enqueue(models.Expression{Expression: "(5-6)/2"})

	expected := []string{"1+2", "3*4", "(5-6)/2"}
	for _, expr := range expected {
		result, ok := q.Dequeue()
		if !ok {
			t.Errorf("Dequeue failed unexpectedly")
		}
		if result.Expression != expr {
			t.Errorf("Dequeue() = %s; want %s", result.Expression, expr)
		}
	}

	_, ok := q.Dequeue()
	if ok {
		t.Errorf("Dequeue should fail on empty queue")
	}
}

func TestEnqueueList(t *testing.T) {
	q := queueMaster.ExpressionQueue()

	data := []models.Expression{
		{Expression: "1+2"},
		{Expression: "3*4"},
		{Expression: "(5-6)/2"},
	}
	q.EnqueueList(data)

	expected := []string{"1+2", "3*4", "(5-6)/2"}
	for _, expr := range expected {
		result, ok := q.Dequeue()
		if !ok {
			t.Errorf("Dequeue failed unexpectedly")
		}
		if result.Expression != expr {
			t.Errorf("Dequeue() = %s; want %s", result.Expression, expr)
		}
	}

	_, ok := q.Dequeue()
	if ok {
		t.Errorf("Dequeue should fail on empty queue")
	}

	q.EnqueueList([]models.Expression{})
	_, ok = q.Dequeue()
	if ok {
		t.Errorf("Dequeue should fail on empty queue after enqueuing an empty list")
	}
}

func TestEmptyDequeue(t *testing.T) {
	q := queueMaster.ExpressionQueue()

	_, ok := q.Dequeue()
	if ok {
		t.Errorf("Dequeue should fail on empty queue")
	}
}
