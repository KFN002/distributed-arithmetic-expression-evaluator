package calculator

import "strings"

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

func InfixToPostfix(expression string) []string {
	var postfixExpression []string
	var operators []string

	expression = strings.ReplaceAll(expression, " ", "")
	number := ""

	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			number += string(char)
		case char == '(':
			handleNumber(&postfixExpression, &number)
			operators = append(operators, "(")
		case char == ')':
			handleNumber(&postfixExpression, &number)
			handleClosingParenthesis(&postfixExpression, &operators)
		default:
			handleNumber(&postfixExpression, &number)
			handleOperator(char, &postfixExpression, &operators)
		}
	}

	handleNumber(&postfixExpression, &number)
	handleRemainingOperators(&postfixExpression, &operators)

	return postfixExpression
}

func handleNumber(postfixExpression *[]string, number *string) {
	if *number != "" {
		*postfixExpression = append(*postfixExpression, *number)
		*number = ""
	}
}

func handleClosingParenthesis(postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 && (*operators)[len(*operators)-1] != "(" {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
	*operators = (*operators)[:len(*operators)-1]
}

func handleOperator(char rune, postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 && priority(rune((*operators)[len(*operators)-1][0])) >= priority(char) {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
	*operators = append(*operators, string(char))
}

func handleRemainingOperators(postfixExpression *[]string, operators *[]string) {
	for len(*operators) > 0 {
		*postfixExpression = append(*postfixExpression, (*operators)[len(*operators)-1])
		*operators = (*operators)[:len(*operators)-1]
	}
}
