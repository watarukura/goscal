package main

import (
	"fmt"
	"log"
	"strconv"
)

type Expression Valuer
type Ident string
type Number float64
type Valuer interface {
	getValue() interface{}
}

func (i Ident) getValue() interface{}  { return i }
func (n Number) getValue() interface{} { return n }

// func (l LParen) getValue() interface{} { return l }
// func (r RParen) getValue() interface{} { return r }
func (a Add) getValue() interface{} { return a }

type Add struct {
	Left, Right Expression
}

func main() {
	testCases := []string{
		//"123",
		//"Hello + world",
		//"(123 + 456 ) + world",
		//"car + cdr + cdr",
		//"((1 + 2) + (3 + 4)) + 5 + 6",
		"1 + 2",
	}
	for _, input := range testCases {
		_, expr, err := expr(input)
		if err != nil {
			log.Fatalf("failed to parse expression '%s': %v", input, err)
		} else {
			fmt.Printf("source: %q, parsed: %#v\n", input, expr)
		}
	}
}

func whitespace(input string) string {
	for len(input) > 0 && input[0] == ' ' {
		input = input[1:]
	}
	return input
}
func number(input string) (string, Expression, error) {
	if len(input) > 0 && (input[0] == '-' || input[0] == '+' || input[0] == '.' || ('0' <= input[0] && input[0] <= '9')) {
		i := 0
		for ; i < len(input) && (input[i] == '.' || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		num, _ := strconv.ParseFloat(input[:i], 64)
		input = input[i:]
		return input, Number(num), nil
	}
	return input, nil, fmt.Errorf("number: invalid input: %q", input)
}
func ident(input string) (string, Expression, error) {
	if len(input) > 0 && (('a' <= input[0] && input[0] <= 'z') || ('A' <= input[0] && input[0] <= 'Z')) {
		i := 0
		for ; i < len(input) && (('a' <= input[i] && input[i] <= 'z') || ('A' <= input[i] && input[i] <= 'Z') || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		ident := input[:i]
		input = input[i:]
		return input, Ident(ident), nil
	}
	return input, nil, fmt.Errorf("ident: invalid input: %q", input)
}

func token(input string) (string, Expression, error) {
	if input, ident, err := ident(whitespace(input)); err == nil {
		return input, ident, nil
	}
	if input, number, err := number(whitespace(input)); err == nil {
		return input, number, nil
	}
	return input, nil, fmt.Errorf("token: invalid input: %q", input)
}

func lparen(input string) string {
	if len(input) > 0 && input[0] == '(' {
		return input[1:]
	}
	return ""
}
func rparen(input string) string {
	if len(input) > 0 && input[0] == ')' {
		return input[1:]
	}
	return ""
}

func plus(input string) string {
	if len(input) > 0 && input[0] == '+' {
		if len(input) > 1 {
			return input[1:]
		}
	}
	return input
}

func expr(input string) (string, Expression, error) {
	if res, expr, err := add(input); err == nil {
		return res, expr, nil
	}

	if res, expr, err := term(input); err == nil {
		return res, expr, nil
	}

	return "", nil, fmt.Errorf("failed to match pattern")
}

func paren(input string) (string, Expression, error) {
	nextInput := lparen(whitespace(input))
	if nextInput == "" {
		return "", nil, fmt.Errorf("failed to match pattern")
	}
	nextInput, expr, err := expr(nextInput)
	nextInput = rparen(whitespace(nextInput))

	return nextInput, expr, err
}

func addTerm(input string) (string, Expression, error) {
	nextInput, lhs, err := term(whitespace(input))
	nextInput = plus(whitespace(nextInput))

	return nextInput, lhs, err
}

func add(input string) (string, Expression, error) {
	var left, right Expression
	//for {
	//	nextInput, expr, err := addTerm(input)
	//	if err != nil {
	//		break
	//	}
	//	if left != nil {
	//		newExpr := Add{left, expr}
	//		left = (Expression)(newExpr)
	//	} else {
	//		left = expr
	//	}
	//	input = nextInput
	//}
	//
	//if left == nil {
	//	return input, nil, fmt.Errorf("failed to match pattern")
	//}

	nextInput, left, err := addTerm(input)
	nextInput, right, err2 := addTerm(nextInput)
	if err == nil && err2 == nil {
		return nextInput, Add{left, right}, err
	}

	return nextInput, nil, fmt.Errorf("failed to match pattern")
}

func term(input string) (string, Expression, error) {
	if res, expr, err := paren(input); err == nil {
		return res, expr, nil
	}

	if res, expr, err := token(input); err == nil {
		return res, expr, nil
	}

	return "", nil, fmt.Errorf("failed to match pattern")
}
