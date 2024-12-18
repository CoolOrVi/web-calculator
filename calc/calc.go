package calc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLeftParen
	TokenRightParen
)

type Token struct {
	Type  TokenType
	Value string
}

func Tokenize(expr string) ([]Token, error) {
	var tokens []Token
	var currentToken strings.Builder

	for i, char := range expr {
		if unicode.IsDigit(char) || char == '.' {
			currentToken.WriteRune(char)
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenNumber, Value: currentToken.String()})
				currentToken.Reset()
			}

			if char == '-' && (i == 0 || tokens[len(tokens)-1].Type == TokenOperator || tokens[len(tokens)-1].Type == TokenLeftParen) {
				currentToken.WriteRune(char)
			} else {
				tokens = append(tokens, Token{Type: TokenOperator, Value: string(char)})
			}
		} else if char == '(' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenNumber, Value: currentToken.String()})
				currentToken.Reset()
			}
			tokens = append(tokens, Token{Type: TokenLeftParen, Value: string(char)})
		} else if char == ')' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, Token{Type: TokenNumber, Value: currentToken.String()})
				currentToken.Reset()
			}
			tokens = append(tokens, Token{Type: TokenRightParen, Value: string(char)})
		} else if unicode.IsSpace(char) {
			continue
		} else {
			return nil, fmt.Errorf("invalid character: %c", char)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, Token{Type: TokenNumber, Value: currentToken.String()})
	}

	return tokens, nil
}

func ParseTokens(tokens []Token) (float64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}

	var values []float64
	var operators []string

	applyOperator := func() error {
		if len(values) < 2 {
			return fmt.Errorf("error in expression")
		}

		b := values[len(values)-1]
		a := values[len(values)-2]
		values = values[:len(values)-2]

		operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]

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
				return fmt.Errorf("division by zero")
			}
			result = a / b
		default:
			return fmt.Errorf("unknown operator: %s", operator)
		}

		values = append(values, result)
		return nil
	}

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("incorrect number: %s", token.Value)
			}
			values = append(values, value)
		case TokenOperator:
			for len(operators) > 0 && operators[len(operators)-1] != "(" &&
				precedence[operators[len(operators)-1]] >= precedence[token.Value] {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token.Value)
		case TokenLeftParen:
			operators = append(operators, token.Value)
		case TokenRightParen:
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 || operators[len(operators)-1] != "(" {
				return 0, fmt.Errorf("bracket mismatch")
			}
			operators = operators[:len(operators)-1]
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return 0, fmt.Errorf("bracket mismatch")
		}
		if err := applyOperator(); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, fmt.Errorf("error in expression")
	}

	return values[0], nil
}
