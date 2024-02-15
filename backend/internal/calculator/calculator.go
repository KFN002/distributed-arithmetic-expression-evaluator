package calculator

import (
	"strconv"
	"time"
)

func Solve(tokens []string, operations map[string]int) int {
	var stack []int
	for _, elem := range tokens {
		if elem == "+" || elem == "-" || elem == "*" || elem == "/" {
			firstNum := stack[len(stack)-2]
			secondNum := stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			time.Sleep(time.Second * time.Duration(operations[elem]))
			switch elem {
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
			num, _ := strconv.Atoi(elem)
			stack = append(stack, num)
		}
	}
	return stack[0]
}
