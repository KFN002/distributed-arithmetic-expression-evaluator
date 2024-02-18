package utils

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"regexp"
	"strings"
	"unicode"
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
	if !hasNoProblemsWithOperators(expression) {
		return false
	}
	if !operatorsFromEachSide(expression) {
		return false
	}
	// Если все проверки пройдены, возвращаем true
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

// Проверяет, что операторы +, -, *, / не стоят в начале или в конце строки
func hasNoProblemsWithOperators(expression string) bool {
	// Убираем пробелы из выражения
	expression = strings.ReplaceAll(expression, " ", "")
	if len(expression) == 0 {
		return false // Пустая строка не является допустимым выражением
	}

	// Проверяем, что первый и последний символ не оператор
	if isOperator(rune(expression[0])) || isOperator(rune(expression[len(expression)-1])) {
		return false
	}

	// Проверяем наличие унарных минусов в начале выражения или после открывающей скобки
	for i, char := range expression {
		if char == '-' && (i == 0 || expression[i-1] == '(') {
			return false
		}
	}

	// Проверяем наличие унарных минусов после закрывающей скобки
	for i := 0; i < len(expression)-1; i++ {
		if expression[i] == ')' && expression[i+1] == '-' {
			return false
		}
	}

	return true
}

func operatorsFromEachSide(expression string) bool {
	for i := 0; i < len(expression); i++ {
		if expression[i] == '(' {
			for j := i + 1; j < len(expression); j++ {
				if expression[j] == ')' {
					break
				}
				i++
			}
			continue
		}
		if expression[i] == ')' {
			for j := i - 1; j >= 0; j-- {
				if isOperator(rune(expression[j])) {
					break
				}
				if expression[j] != ' ' {
					return false
				}
			}
			continue
		}
		if isOperator(rune(expression[i])) {
			if i == 0 || i == len(expression)-1 {
				return false
			}
			if !unicode.IsDigit(rune(expression[i-1])) && expression[i-1] != ')' {
				return false
			}
			if !unicode.IsDigit(rune(expression[i+1])) && expression[i+1] != '(' {
				return false
			}
			if expression[i-1] == ')' && expression[i+1] == '(' {
				return false
			}
		}
	}
	return true
}

// isOperator проверяет, является ли символ оператором +, -, *, /
func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

// FlipList переворачивает список с данными - reverse
func FlipList(list []models.Expression) []models.Expression {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}
