package calculator

import "strings"

// priority определяет приоритет оператора
func priority(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

// InfixToPostfix конвертирует инфиксное выражение в постфиксное
func InfixToPostfix(expression string) []string {
	var postfixExpression []string // постфиксное выражение
	var operators []string         // стек операторов

	expression = strings.ReplaceAll(expression, " ", "") // убираем пробелы
	number := ""

	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			number += string(char)
		case char == '(':
			handleNumber(&postfixExpression, &number) // добавляем число в выражение
			operators = append(operators, "(")        // добавляем открывающую скобку в стек
		case char == ')':
			handleNumber(&postfixExpression, &number)                // добавляем число в выражение
			handleClosingParenthesis(&postfixExpression, &operators) // обрабатываем закрывающую скобку
		default:
			handleNumber(&postfixExpression, &number)            // добавляем число в выражение
			handleOperator(char, &postfixExpression, &operators) // обрабатываем оператор
		}
	}

	handleNumber(&postfixExpression, &number)                // добавляем оставшееся число в выражение
	handleRemainingOperators(&postfixExpression, &operators) // обрабатываем оставшиеся операторы

	return postfixExpression
}

// handleNumber добавляет число в постфиксное выражение, если оно есть
func handleNumber(postfixExpression *[]string, number *string) {
	if *number != "" {
		*postfixExpression = append(*postfixExpression, *number)
		*number = ""
	}
}

// handleClosingParenthesis обрабатывает закрывающую скобку
func handleClosingParenthesis(postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 && (*operators)[len(*operators)-1] != "(" {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
	*operators = (*operators)[:len(*operators)-1] // удаляем открывающую скобку из стека
}

// handleOperator обрабатывает оператор
func handleOperator(char rune, postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 && priority(rune((*operators)[len(*operators)-1][0])) >= priority(char) {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
	*operators = append(*operators, string(char)) // добавляем текущий оператор в стек
}

// handleRemainingOperators обрабатывает оставшиеся операторы
func handleRemainingOperators(postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
}
