package tests

import (
	"errors"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/calculator"
	"reflect"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		number1   float64
		number2   float64
		operation string
		expected  float64
		err       error
	}{
		{5, 2, "+", 7, nil},
		{5, 2, "-", 3, nil},
		{5, 2, "*", 10, nil},
		{5, 2, "/", 2.5, nil},
		{5, 0, "/", 0, errors.New("division by zero")},
		{5, 2, "%", 0, errors.New("invalid operation")},
	}

	for _, test := range tests {
		result, err := calculator.Calculate(test.number1, test.number2, test.operation)
		if err != nil && err.Error() != test.err.Error() {
			t.Errorf("Calculate(%f, %f, %s) returned unexpected error: got %v, want %v", test.number1, test.number2, test.operation, err, test.err)
		}
		if result != test.expected {
			t.Errorf("Calculate(%f, %f, %s) returned unexpected result: got %f, want %f", test.number1, test.number2, test.operation, result, test.expected)
		}
	}
}

func TestInfixToPostfix(t *testing.T) {
	tests := []struct {
		expression string
		expected   []string
	}{
		{"1+2", []string{"1", "2", "+"}},
		{"(1+2)*3", []string{"1", "2", "+", "3", "*"}},
		{"1*2+3", []string{"1", "2", "*", "3", "+"}},
		{"(1+2)*3+4", []string{"1", "2", "+", "3", "*", "4", "+"}},
	}

	for _, test := range tests {
		result := calculator.InfixToPostfix(test.expression)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("InfixToPostfix(%s) returned unexpected result: got %v, want %v", test.expression, result, test.expected)
		}
	}
}
