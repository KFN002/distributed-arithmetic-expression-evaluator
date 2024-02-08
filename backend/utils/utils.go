package utils

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"strings"
)

// Проверяет выражение на сбалансированность скобок и на отсутствие двух или более арифметических знаков рядом.
func CheckExpression(expression string) bool {
	if !areParenthesesBalanced(expression) {
		return false
	}
	// Проверяем отсутствие двух или более арифметических знаков рядом
	if hasConsecutiveOperators(expression) {
		return false
	}
	// Если обе проверки пройдены, возвращаем true
	return true
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

func FlipList(list []models.Expression) []models.Expression {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// Польская нотация и перевод в нее
func isOperator(ch string) bool {
	operators := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
	}
	return operators[ch]
}

func precedence(ch string) int {
	switch ch {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func PolishNotation(infix string) string {
	var result string
	var stack []string

	tokens := strings.Fields(infix)

	for _, token := range tokens {
		if isOperator(token) {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				result += stack[len(stack)-1] + " "
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {
			result += token + " "
		}
	}

	for len(stack) > 0 {
		result += stack[len(stack)-1] + " "
		stack = stack[:len(stack)-1]
	}

	return strings.TrimSpace(result)
}
