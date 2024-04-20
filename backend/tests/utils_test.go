package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/utils"
	"testing"
)

func TestSumList(t *testing.T) {
	tests := []struct {
		data     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, 6},
		{[]float64{-1, -2, -3}, -6},
		{[]float64{}, 0},
	}

	for _, test := range tests {
		result := utils.SumList(test.data)
		if result != test.expected {
			t.Errorf("SumList(%v) = %f; want %f", test.data, result, test.expected)
		}
	}
}

// Проверка после регулярного выражения
func TestCheckExpression(t *testing.T) {
	tests := []struct {
		expression string
		expected   bool
	}{
		{"(2+3)*5", true},
		{"2+3+5", true},
		{"2+3/1", true},
		{"2*3-1", true},
		{"(2+3*5)", true},
		{"(2+3)*(5-1)", true},
		{"(2+3++5)", false},
		{"(2+(3++5))", false},
		{"2+3+5)", false},
		{"2+3+(5", false},
		{"2+3/0", false},
		{"2+3/", false},
		{"2++3", false},
		{"2+*3", false},
		{"(2+3", false},
		{"2+3*5)", false},
		{"2+(3++5)", false},
		{"+2+3", false},
		{"2+3+", false},
	}

	for _, test := range tests {
		result := utils.CheckExpression(test.expression)
		if result != test.expected {
			t.Errorf("CheckExpression(%s) = %t; want %t", test.expression, result, test.expected)
		}
	}
}

func TestFlipList(t *testing.T) {
	tests := []struct {
		list     []models.Expression
		expected []models.Expression
	}{
		{[]models.Expression{{Expression: "a"}, {Expression: "b"}, {Expression: "c"}}, []models.Expression{{Expression: "c"}, {Expression: "b"}, {Expression: "a"}}},
		{[]models.Expression{}, []models.Expression{}},
	}

	for _, test := range tests {
		result := utils.FlipList(test.list)
		for i := range result {
			if result[i].Expression != test.expected[i].Expression {
				t.Errorf("FlipList(%v) = %v; want %v", test.list, result, test.expected)
				break
			}
		}
	}
}
