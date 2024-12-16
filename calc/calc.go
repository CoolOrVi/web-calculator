package calc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

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

func applyOperator(values []float64, operator string) ([]float64, error) {
	if len(values) < 2 {
		return values, fmt.Errorf("invalid operands")
	}
	b := values[len(values)-1]
	a := values[len(values)-2]
	values = values[:len(values)-2]

	var result float64
	switch operator {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return values, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return values, fmt.Errorf("invalid operands")
	}

	values = append(values, result)
	return values, nil
}

func Calc(expression string) (float64, error) {
	var values []float64
	var operators []string

	expression = strings.ReplaceAll(expression, " ", "")

	var i int
	for i < len(expression) {
		ch := expression[i]

		if unicode.IsDigit(rune(ch)) || ch == '.' {
			j := i
			for j < len(expression) && (unicode.IsDigit(rune(expression[j])) || expression[j] == '.') {
				j++
			}
			numberStr := expression[i:j]
			number, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number")
			}
			values = append(values, number)
			i = j - 1

		} else if ch == '(' {
			operators = append(operators, string(ch))

		} else if ch == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				var err error
				operators, values, err = popAndApplyOperator(operators, values)
				if err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, fmt.Errorf("parenthesis mismatch")
			}
			operators = operators[:len(operators)-1]

		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			currOp := string(ch)
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(currOp) {
				var err error
				operators, values, err = popAndApplyOperator(operators, values)
				if err != nil {
					return 0, err
				}
			}
			operators = append(operators, currOp)

		} else {
			return 0, fmt.Errorf("invalid symbol")
		}
		i++
	}

	for len(operators) > 0 {
		var err error
		operators, values, err = popAndApplyOperator(operators, values)
		if err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, fmt.Errorf("error in expression")
	}

	return values[0], nil
}

func popAndApplyOperator(operators []string, values []float64) ([]string, []float64, error) {
	if len(operators) == 0 {
		return operators, values, fmt.Errorf("invalid operands")
	}
	operator := operators[len(operators)-1]
	operators = operators[:len(operators)-1]

	values, err := applyOperator(values, operator)
	if err != nil {
		return operators, values, err
	}
	return operators, values, nil
}
