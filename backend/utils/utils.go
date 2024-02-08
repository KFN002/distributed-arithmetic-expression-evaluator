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

// Постфиксная нотация
var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func isOperator(ch rune) bool {
	_, ok := precedence[ch]
	return ok
}

func InfixToPostfix(infix string) string {
	var result strings.Builder
	var stack []rune

	for _, token := range infix {
		switch {
		case token >= '0' && token <= '9':
			result.WriteRune(token)
			result.WriteRune(' ')
		case isOperator(token):
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				result.WriteRune(stack[len(stack)-1])
				result.WriteRune(' ')
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case token == '(':
			stack = append(stack, token)
		case token == ')':
			for stack[len(stack)-1] != '(' {
				result.WriteRune(stack[len(stack)-1])
				result.WriteRune(' ')
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		}
	}

	for len(stack) > 0 {
		result.WriteRune(stack[len(stack)-1])
		result.WriteRune(' ')
		stack = stack[:len(stack)-1]
	}

	return strings.TrimSpace(result.String())
}
