package utils

import (
	"distributed-arithmetic-expression-evaluator/backend/pkg/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func SumList(data []float64) float64 {
	var total float64
	for _, elem := range data {
		total += elem
	}
	return total
}

// CheckExpression Проверяет выражение на сбалансированность скобок и на отсутствие двух или более арифметических знаков рядом.
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

func containsOperator(input string) bool {
	operatorRegex := regexp.MustCompile(`[+\-*\/]`)

	return operatorRegex.MatchString(input)
}

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

func RemoveRedundantParentheses(expression string) string {
	var result strings.Builder
	stack := make([]rune, 0)

	for _, char := range expression {
		if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			if len(stack) > 0 && stack[len(stack)-1] == '(' {
				stack = stack[:len(stack)-1]
			} else {
				result.WriteRune(char)
			}
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func FlipList(list []models.Expression) []models.Expression {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func CalculateSimpleTask(expression string) (float64, error) {
	operators := []string{"+", "-", "*", "/"}
	var operator string
	for _, op := range operators {
		if strings.Contains(expression, op) {
			operator = op
			break
		}
	}

	if operator == "" {
		return 0, fmt.Errorf("no valid operator found in the expression")
	}

	parts := strings.Split(expression, operator)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid expression format")
	}

	left, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid left operand: %v", err)
	}

	right, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid right operand: %v", err)
	}

	var result float64
	switch operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero error")
		}
		result = left / right
	}

	return result, nil
}

func CountOperators(inputString string) int {
	operators := []string{"+", "-", "*", "/"}
	operatorCount := 0

	for _, op := range operators {
		operatorCount += strings.Count(inputString, op)
	}

	return operatorCount
}
