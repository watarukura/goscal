package main

import (
	"fmt"
	"log"
	"math"
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
		"pi",
		"1",
		"1 + 2 + 3",
		"(123 + 456 ) + pi",
		"10 + (100 + 1)",
		"((1 + 2) + (3 + 4)) + 5 + 6",
	}
	for _, input := range testCases {
		_, expr, _ := expr(input)
		evaled := eval(expr)
		fmt.Printf("source: %q, parsed: %#v\n", input, evaled)
	}
}

func eval(expr Expression) Number {
	switch expr.(type) {
	case Ident:
		if expr == Ident("pi") {
			return math.Pi
		}
		log.Fatalf("Unknown identifier '%v'", expr)
	case Number:
		return expr.(Number)
	case Add:
		return eval(expr.(Add).Left) + eval(expr.(Add).Right)
	}
	return Number(0)
}

func whitespace(input string) string {
	for len(input) > 0 && input[0] == ' ' {
		input = input[1:]
	}
	return input
}
func number(input string) (string, Expression, bool) {
	if len(input) > 0 && (input[0] == '-' || input[0] == '+' || input[0] == '.' || ('0' <= input[0] && input[0] <= '9')) {
		i := 0
		for ; i < len(input) && (input[i] == '.' || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		num, _ := strconv.ParseFloat(input[:i], 64)
		input = input[i:]
		return input, Number(num), true
	}
	return input, nil, false
}
func ident(input string) (string, Expression, bool) {
	if len(input) > 0 && (('a' <= input[0] && input[0] <= 'z') || ('A' <= input[0] && input[0] <= 'Z')) {
		i := 0
		for ; i < len(input) && (('a' <= input[i] && input[i] <= 'z') || ('A' <= input[i] && input[i] <= 'Z') || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		ident := input[:i]
		input = input[i:]
		return input, Ident(ident), true
	}
	return input, nil, false
}

func token(input string) (string, Expression, bool) {
	if input, ident, ok := ident(whitespace(input)); ok {
		return input, ident, true
	}
	if input, number, ok := number(whitespace(input)); ok {
		return input, number, true
	}
	return input, nil, false
}

func plus(input string) string {
	if len(input) > 0 && input[0] == '+' {
		if len(input) > 1 {
			return input[1:]
		}
	}
	return input
}

func expr(input string) (string, Expression, bool) {
	if res, expr, ok := add(input); ok {
		return res, expr, true
	}

	if res, expr, ok := term(input); ok {
		return res, expr, true
	}

	return "", nil, false
}

func paren(input string) (string, Expression, bool) {
	input = whitespace(input)
	nextInput, ok := lparen(input)

	if ok {
		nextInput, expr, ok := expr(nextInput)
		nextInput, _ = rparen(whitespace(nextInput))
		return nextInput, expr, ok
	}

	nextInput, ok = rparen(input)
	if ok {
		return nextInput, nil, true
	}

	return input, nil, false
}
func lparen(input string) (string, bool) {
	if len(input) > 0 && input[0] == '(' {
		return input[1:], true
	}
	return input, false
}
func rparen(input string) (string, bool) {
	if len(input) > 0 && input[0] == ')' {
		return input[1:], true
	}
	return input, false
}

func addTerm(input string) (string, Expression, bool) {
	nextInput, lhs, ok := term(whitespace(input))
	nextInput = plus(whitespace(nextInput))
	if ok {
		return nextInput, lhs, true
	}
	return input, nil, false
}

func add(input string) (string, Expression, bool) {
	var add Expression
	var nextInput string
	for {
		nextInput, expr, ok := addTerm(input)
		if !ok {
			break
		}
		if expr == nil {
			input = nextInput
			continue
		}
		if add != nil {
			add = Add{add, expr}
		} else {
			add = expr
		}
		input = nextInput
	}

	if add == nil {
		return nextInput, nil, false
	}
	return nextInput, add, true
}

func term(input string) (string, Expression, bool) {
	if res, expr, ok := paren(input); ok {
		return res, expr, true
	}

	if res, expr, ok := token(input); ok {
		return res, expr, true
	}

	return "", nil, false
}
