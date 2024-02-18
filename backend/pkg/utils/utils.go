package utils

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"regexp"
	"strings"
)

// CheckExpression проверяет выражение на сбалансированность скобок и на отсутствие двух или более арифметических знаков рядом.
func CheckExpression(expression string) bool {
	if !areParenthesesBalanced(expression) {
		return false
	}
	// Проверяем отсутствие двух или более арифметических знаков рядом
	if hasConsecutiveOperators(expression) {
		return false
	}
	if !containsOperator(expression) {
		return false
	}
	if HasDivisionByZero(expression) {
		return false
	}
	// Если обе проверки пройдены, возвращаем true
	return true
}

// Проверка на наличие арифметического оператора
func containsOperator(input string) bool {
	operatorRegex := regexp.MustCompile(`[+\-*\/]`)

	return operatorRegex.MatchString(input)
}

// HasDivisionByZero проверка на выраженное деление на 0
func HasDivisionByZero(expression string) bool {
	operands := strings.Split(expression, "/")

	for _, op := range operands {
		if op == "0" {
			return true
		}
	}
	return false
}

// Проверяет, сбалансированы ли скобки в выражении
func areParenthesesBalanced(expression string) bool {
	stack := make([]rune, 0)

	for _, char := range expression {
		if char == '(' {
			stack = append(stack, '(')
		} else if char == ')' {
			if len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

// Проверяет, есть ли в выражении два или более арифметических знаков рядом
func hasConsecutiveOperators(expression string) bool {
	operators := "+-*/"
	for i := 0; i < len(expression)-1; i++ {
		if strings.ContainsAny(string(expression[i]), operators) && strings.ContainsAny(string(expression[i+1]), operators) {
			return true
		}
	}
	return false
}

// FlipList переворачивает список с данными - reverse
func FlipList(list []models.Expression) []models.Expression {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}
