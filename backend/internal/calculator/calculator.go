package calculator

import "errors"

func Calculate(firstNum, secondNum float64, operation string) (float64, error) {
	switch operation {
	case "+":
		return firstNum + secondNum, nil
	case "-":
		return firstNum - secondNum, nil
	case "*":
		return firstNum * secondNum, nil
	case "/":
		if secondNum == 0 {
			return 0, errors.New("division by zero")
		}
		return firstNum / secondNum, nil
	default:
		return 0, errors.New("invalid operation")
	}
}
