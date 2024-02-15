package calculator

import (
	"strconv"
	"time"
)

func Solve(tokens []string, operations map[string]int) int {
	var stack []int
	for _, el := range tokens {
		if el == "+" || el == "-" || el == "*" || el == "/" {
			firstNum := stack[len(stack)-2]
			secondNum := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			time.Sleep(time.Second * time.Duration(operations[el]))
			switch el {
			case "+":
				stack = append(stack, firstNum+secondNum)
			case "-":
				stack = append(stack, firstNum-secondNum)
			case "*":
				stack = append(stack, firstNum*secondNum)
			case "/":
				stack = append(stack, firstNum/secondNum)
			}
		} else {
			num, _ := strconv.Atoi(el)
			stack = append(stack, num)
		}
	}
	return stack[0]
}
