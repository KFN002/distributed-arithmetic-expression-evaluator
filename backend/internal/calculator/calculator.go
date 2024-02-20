package calculator

import "errors"

// Calculate Подсчет простого выражения с обработкой случаев с делением на ноль, другими ошибками
func Calculate(number1 float64, number2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		return number1 + number2, nil
	case "-":
		return number1 - number2, nil
	case "*":
		return number1 * number2, nil
	case "/":
		if number2 == 0 {
			return 0, errors.New("division by zero")
		}
		return number1 / number2, nil
	default:
		return 0, errors.New("invalid operation")
	}
}
