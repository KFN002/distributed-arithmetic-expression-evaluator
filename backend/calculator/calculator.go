package calculator

import (
	"strconv"
	"strings"
	"time"
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
	var stack []string

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
			stack = append(stack, "(")
		case char == ')':
			if number != "" {
				result = append(result, number)
				number = ""
			}
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
		default:
			if number != "" {
				result = append(result, number)
				number = ""
			}
			for len(stack) > 0 && priority(rune(stack[len(stack)-1][0])) >= priority(char) {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, string(char))
		}
	}

	if number != "" {
		result = append(result, number)
	}

	for len(stack) > 0 {
		result = append(result, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return result
}

func Calculate(tokens []string, operations map[string]int) int {
	var stack []int
	for _, el := range tokens {
		if el == "+" || el == "-" || el == "*" || el == "/" {
			firstNum := stack[len(stack)-2]
			secondNum := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			if el == "-" {
				time.Sleep(time.Second * time.Duration(operations[el]))
				stack = append(stack, firstNum-secondNum)
			} else if el == "+" {
				time.Sleep(time.Second * time.Duration(operations[el]))
				stack = append(stack, firstNum+secondNum)
			} else if el == "*" {
				time.Sleep(time.Second * time.Duration(operations[el]))
				stack = append(stack, firstNum*secondNum)
			} else {
				time.Sleep(time.Second * time.Duration(operations[el]))
				stack = append(stack, firstNum/secondNum)
			}
		} else {
			num, _ := strconv.Atoi(el)
			stack = append(stack, num)
		}
	}
	return stack[0]
}
