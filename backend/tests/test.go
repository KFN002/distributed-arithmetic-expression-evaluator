package main

import (
	"fmt"
)

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

func main() {
	expressions := []string{"(a+b)-(c-d)", "a+b(c-d)", "a+(b*c)-d"}

	for _, exp := range expressions {
		if hasOperatorNearParentheses(exp) {
			fmt.Printf("String '%s' has operators near parentheses where required.\n", exp)
		} else {
			fmt.Printf("String '%s' doesn't have operators near parentheses where required.\n", exp)
		}
	}
}
