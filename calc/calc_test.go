package calc_test

import (
	"testing"

	"github.com/coolorvi/web-calculator/calc"

)

func TestTokenize(t *testing.T) {
	tests := []struct {
		expr     string
		expected []Token
		isError  bool
	}{
		{"3 + 5", []Token{
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
		}, false},
		{"-3 + 5", []Token{
			{Type: TokenNumber, Value: "-3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
		}, false},
		{"(3 + 5) * 2", []Token{
			{Type: TokenLeftParen, Value: "("},
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
			{Type: TokenRightParen, Value: ")"},
			{Type: TokenOperator, Value: "*"},
			{Type: TokenNumber, Value: "2"},
		}, false},
		{"3.5 * -2", []Token{
			{Type: TokenNumber, Value: "3.5"},
			{Type: TokenOperator, Value: "*"},
			{Type: TokenNumber, Value: "-2"},
		}, false},
		{"3 + #", nil, true},
	}

	for _, test := range tests {
		tokens, err := tokenize(test.expr)
		if test.isError {
			if err == nil {
				t.Errorf("expected error for expression: %s", test.expr)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for expression %s: %v", test.expr, err)
			continue
		}
		if len(tokens) != len(test.expected) {
			t.Errorf("unexpected token count for expression %s: got %d, want %d", test.expr, len(tokens), len(test.expected))
			continue
		}
		for i, token := range tokens {
			if token != test.expected[i] {
				t.Errorf("unexpected token at position %d for expression %s: got %+v, want %+v", i, test.expr, token, test.expected[i])
			}
		}
	}
}

func TestParseTokens(t *testing.T) {
	tests := []struct {
		tokens   []Token
		expected float64
		isError  bool
	}{
		{[]Token{
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
		}, 8, false},
		{[]Token{
			{Type: TokenNumber, Value: "-3"},
			{Type: TokenOperator, Value: "*"},
			{Type: TokenNumber, Value: "5"},
		}, -15, false},
		{[]Token{
			{Type: TokenLeftParen, Value: "("},
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
			{Type: TokenRightParen, Value: ")"},
			{Type: TokenOperator, Value: "*"},
			{Type: TokenNumber, Value: "2"},
		}, 16, false},
		{[]Token{
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "/"},
			{Type: TokenNumber, Value: "0"},
		}, 0, true},
		{[]Token{
			{Type: TokenLeftParen, Value: "("},
			{Type: TokenNumber, Value: "3"},
			{Type: TokenOperator, Value: "+"},
			{Type: TokenNumber, Value: "5"},
		}, 0, true},
	}

	for _, test := range tests {
		result, err := Calc(test.tokens)
		if test.isError {
			if err == nil {
				t.Errorf("expected error for tokens: %+v", test.tokens)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for tokens %+v: %v", test.tokens, err)
			continue
		}
		if result != test.expected {
			t.Errorf("unexpected result for tokens %+v: got %f, want %f", test.tokens, result, test.expected)
		}
	}
}

func TestCalculatorIntegration(t *testing.T) {
	tests := []struct {
		expr     string
		expected float64
		isError  bool
	}{
		{"3 + 5", 8, false},
		{"-3 + 5", 2, false},
		{"(3 + 5) * 2", 16, false},
		{"3.5 * -2", -7, false},
		{"(3 + 5) * (10 - 2) / 2", 32, false},
		{"3 + 5 / (2 - 2)", 0, true},
		{"3 + 5 *", 0, true},
	}

	for _, test := range tests {
		tokens, err := tokenize(test.expr)
		if err != nil {
			if !test.isError {
				t.Errorf("unexpected error during tokenization for expression %s: %v", test.expr, err)
			}
			continue
		}

		result, err := Calc(tokens)
		if test.isError {
			if err == nil {
				t.Errorf("expected error for expression: %s", test.expr)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for expression %s: %v", test.expr, err)
			continue
		}
		if result != test.expected {
			t.Errorf("unexpected result for expression %s: got %f, want %f", test.expr, result, test.expected)
		}
	}
}
