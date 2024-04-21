package utils

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"regexp"
	"strings"
	"unicode"
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
	if !checkPrefixSuffix(expression) {
		return false
	}
	if !areParenthesesBalanced(expression) {
		return false
	}
	if !containsOperator(expression) {
		return false
	}
	if hasConsecutiveOperators(expression) {
		return false
	}
	if hasDivisionByZero(expression) {
		return false
	}
	if !hasOperatorNearParentheses(expression) {
		return false
	}
	// Если обе проверки пройдены, возвращаем true
	return true
}

func containsOperator(input string) bool {
	operatorRegex := regexp.MustCompile(`[+\-*\/]`)

	return operatorRegex.MatchString(input)
}

func hasDivisionByZero(expression string) bool {
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

// Проверяет правильность операторов у скобок - не используется с 1.02
func isValidExpression(expr string) bool {

	if strings.HasPrefix(expr, "(") || strings.HasSuffix(expr, ")") {
		return false
	}

	if strings.Contains(expr, "(") {
		if !unicode.IsDigit(rune(expr[0])) || !unicode.IsDigit(rune(expr[len(expr)-1])) {
			return false
		}
	}

	parts := strings.Split(expr, "(")
	for _, part := range parts {
		if strings.Contains(part, ")") {
			if strings.ContainsAny(part, "+-*/") {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

// Проверяет правильность расстановки операторов в выражении
func checkPrefixSuffix(expression string) bool {
	var operators = []string{"+", "-", "/", "*"}
	for _, operator := range operators {
		if strings.HasPrefix(expression, operator) || strings.HasSuffix(expression, operator) {
			return false
		}
	}
	return true
}

// Проверка правильности знаков у скобок
func hasOperatorNearParentheses(expression string) bool {
	for i := 1; i < len(expression)-1; i++ {
		if expression[i] == '(' {
			if !isOperator(expression[i-1]) && !isOperator(expression[i+1]) {
				return false
			}
		} else if expression[i] == ')' {
			if !isOperator(expression[i-1]) && !isOperator(expression[i+1]) {
				return false
			}
		}
	}
	return true
}

func isOperator(char byte) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// FlipList переворачивает список с данными - reverse
func FlipList(list []models.Expression) []models.Expression {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}
