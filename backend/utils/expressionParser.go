package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Node представляет узел в дереве выражения.
type Node struct {
	Left     *Node   // Левый дочерний узел
	Right    *Node   // Правый дочерний узел
	Operator string  // Оператор (+, -, *, /)
	Value    float64 // Значение (если это листовой узел)
}

// ExpressionParser разбирает строку входного выражения в дерево выражения.
func ExpressionParser(s string) (*Node, error) {
	var (
		tokens    = strings.Fields(s) // Токенизация входного выражения
		stack     []*Node             // Стек для хранения узлов
		operators []string            // Стек для хранения операторов
	)

	// Итерация по каждому токену в выражении
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			// Если токен - оператор, обрабатываем приоритет
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				popOperator(&stack, &operators)
			}
			operators = append(operators, token)
		default:
			// Если токен не оператор, это число
			value, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			stack = append(stack, &Node{Value: value})
		}
	}

	// Извлечение оставшихся операторов и построение дерева выражения
	for len(operators) > 0 {
		popOperator(&stack, &operators)
	}

	// Проверка корректности дерева выражения
	if len(stack) != 1 {
		return nil, errors.New("неверное выражение")
	}

	return stack[0], nil
}

// precedence назначает уровни приоритета операторам.
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// popOperator извлекает операторы и создает узлы в дереве выражения.
func popOperator(stack *[]*Node, operators *[]string) {
	operator := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	right := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	left := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	node := &Node{Right: right, Left: left, Operator: operator}
	*stack = append(*stack, node)
}

// EvaluatePostOrder вычисляет дерево выражения и сохраняет подвыражения.
func EvaluatePostOrder(node *Node, subExpressions *map[int]string, counter *int) error {
	if node == nil {
		return nil
	}

	if node.Left != nil {
		err := EvaluatePostOrder(node.Left, subExpressions, counter)
		if err != nil {
			return err
		}
	}

	if node.Right != nil {
		err := EvaluatePostOrder(node.Right, subExpressions, counter)
		if err != nil {
			return err
		}
	}

	if node.Left == nil && node.Right == nil {
		(*subExpressions)[*counter] = fmt.Sprintf("%.2f", node.Value)
		*counter++
	}

	if node.Operator != "" {
		lastIndex := *counter - 1
		secondLastIndex := lastIndex - 1
		subExpression := fmt.Sprintf("%s %s %s", (*subExpressions)[secondLastIndex], node.Operator, (*subExpressions)[lastIndex])
		(*subExpressions)[*counter] = subExpression
		*counter++
	}
	return nil
}

// ValidatedPostOrder разбирает выражение и возвращает подвыражения.
func ValidatedPostOrder(s string) (map[int]string, error) {
	node, err := ExpressionParser(s)
	if err != nil {
		return nil, err
	}
	subExps := make(map[int]string)
	var counter int
	err = EvaluatePostOrder(node, &subExps, &counter)
	if err != nil {
		return nil, err
	}
	for key, val := range subExps {
		if len(val) == 4 {
			delete(subExps, key)
		}
	}
	return subExps, nil
}

// Exprt - это пример функции, демонстрирующей использование ValidatedPostOrder.
func Exprt() {
	expression := "2 + 2 * 2"
	subExpressions, err := ValidatedPostOrder(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	// Вывод подвыражений
	for key, val := range subExpressions {
		fmt.Printf("%d: %s\n", key, val)
	}
}
