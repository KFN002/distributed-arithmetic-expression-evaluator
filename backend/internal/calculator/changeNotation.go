package calculator

import (
	"strings"
)

func priority(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func InfixToPostfix(expression string) []string {
	var result []string
	var slice []string

	expression = strings.ReplaceAll(expression, " ", "")
	number := ""

	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			number += string(char)
		case char == '(':
			if number != "" {
				result = append(result, number)
				number = ""
			}
			slice = append(slice, "(")
		case char == ')':
			if number != "" {
				result = append(result, number)
				number = ""
			}
			for len(slice) > 0 && slice[len(slice)-1] != "(" {
				result = append(result, slice[len(slice)-1])
				slice = slice[:len(slice)-1]
			}
			slice = slice[:len(slice)-1]
		default:
			if number != "" {
				result = append(result, number)
				number = ""
			}
			for len(slice) > 0 && priority(rune(slice[len(slice)-1][0])) >= priority(char) {
				result = append(result, slice[len(slice)-1])
				slice = slice[:len(slice)-1]
			}
			slice = append(slice, string(char))
		}
	}

	if number != "" {
		result = append(result, number)
	}

	for len(slice) > 0 {
		result = append(result, slice[len(slice)-1])
		slice = slice[:len(slice)-1]
	}

	return result
}
