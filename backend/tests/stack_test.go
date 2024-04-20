package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"testing"
)

func TestStack_IsEmpty(t *testing.T) {
	stack := models.Stack{}

	if !stack.IsEmpty() {
		t.Errorf("IsEmpty() failed for empty stack. Expected: true")
	}

	stack.Push("item")

	if stack.IsEmpty() {
		t.Errorf("IsEmpty() failed for non-empty stack. Expected: false")
	}
}

func TestStack_Push(t *testing.T) {
	stack := models.Stack{}

	stack.Push("item")

	if len(stack) != 1 {
		t.Errorf("Push() failed. Stack length should be 1 after pushing, got %d", len(stack))
	}

	if stack[0] != "item" {
		t.Errorf("Push() failed. Top element should be 'item', got %s", stack[0])
	}
}

func TestStack_Pop(t *testing.T) {
	stack := models.Stack{"item1", "item2", "item3"}

	popped := stack.Pop()

	if !popped {
		t.Errorf("Pop() failed. Should return true for non-empty stack")
	}

	if len(stack) != 2 {
		t.Errorf("Pop() failed. Stack length should be 2 after popping, got %d", len(stack))
	}

	if stack.Top() != "item2" {
		t.Errorf("Pop() failed. Top element should be 'item2', got %s", stack.Top())
	}

	emptyStack := models.Stack{}
	poppedEmpty := emptyStack.Pop()

	if poppedEmpty {
		t.Errorf("Pop() failed. Should return false for empty stack")
	}
}

func TestStack_Top(t *testing.T) {
	stack := models.Stack{"item1", "item2", "item3"}

	top := stack.Top()

	if top != "item3" {
		t.Errorf("Top() failed. Top element should be 'item3', got %s", top)
	}

	emptyStack := models.Stack{}
	topEmpty := emptyStack.Top()

	if topEmpty != "" {
		t.Errorf("Top() failed. Should return empty string for empty stack")
	}
}
