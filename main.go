package main

import "fmt"

type Token interface {
	getValue() interface{}
}
type Ident string
type Number string

func (i Ident) getValue() interface{}  { return i }
func (n Number) getValue() interface{} { return n }

func main() {
	input := "123 world"
	fmt.Printf("source: %#v, parsed: %#v\n", input, source(input))

	input = "Hello world"
	fmt.Printf("source: %#v, parsed: %#v\n", input, source(input))

	input = "      world"
	fmt.Printf("source: %#v, parsed: %#v\n", input, source(input))
}

func whitespace(input string) string {
	for len(input) > 0 && input[0] == ' ' {
		input = input[1:]
	}
	return input
}
func number(input string) (string, Token) {
	if len(input) > 0 && (input[0] == '-' || input[0] == '+' || input[0] == '.' || ('0' <= input[0] && input[0] <= '9')) {
		i := 0
		for ; i < len(input) && (input[i] == '.' || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		num := input[:i]
		input = input[i:]
		return input, Number(num)
	}
	return input, nil
}
func ident(input string) (string, Token) {
	if len(input) > 0 && (('a' <= input[0] && input[0] <= 'z') || ('A' <= input[0] && input[0] <= 'Z')) {
		i := 0
		for ; i < len(input) && (('a' <= input[i] && input[i] <= 'z') || ('A' <= input[i] && input[i] <= 'Z') || ('0' <= input[i] && input[i] <= '9')); i++ {
		}
		ident := input[:i]
		input = input[i:]
		return input, Ident(ident)
	}
	return input, nil
}

func token(input string) (string, Token) {
	if input, ident := ident(whitespace(input)); ident != nil {
		return input, ident
	}
	if input, number := number(whitespace(input)); number != nil {
		return input, number
	}
	return input, nil
}

func source(input string) []Token {
	var tokens []Token
	var nextToken Token
	for {
		input, nextToken = token(input)
		tokens = append(tokens, nextToken)
		if input == "" {
			break
		}
	}
	return tokens
}
