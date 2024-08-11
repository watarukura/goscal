package main

import "fmt"

type Token interface {
	getValue() interface{}
}
type Ident string
type Number string
type LParen string
type RParen string

func (i Ident) getValue() interface{}  { return i }
func (n Number) getValue() interface{} { return n }
func (l LParen) getValue() interface{} { return l }
func (r RParen) getValue() interface{} { return r }

func main() {
	input := "(123 456 world)"
	fmt.Printf("source: %#v, parsed: %#v\n", input, source(input))

	input = "((car cdr) cdr)"
	fmt.Printf("source: %#v, parsed: %#v\n", input, source(input))

	input = "()())))((()))"
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
	if input, lparen := lparen(whitespace(input)); lparen != nil {
		return input, lparen
	}
	if input, rparen := rparen(whitespace(input)); rparen != nil {
		return input, rparen
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

func lparen(input string) (string, Token) {
	if len(input) > 0 && input[0] == '(' {
		if len(input) > 1 {
			return input[1:], LParen(input[0])
		}
		return "", LParen(input[0])
	}
	return input, nil
}
func rparen(input string) (string, Token) {
	if len(input) > 0 && input[0] == ')' {
		if len(input) > 1 {
			return input[1:], RParen(input[0])
		}
		return "", RParen(input[0])
	}
	return input, nil
}
