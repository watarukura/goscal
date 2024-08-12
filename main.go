package main

import "fmt"

type Token Valuer
type Ident string
type Number string
type LParen string
type RParen string

type TokenTree struct {
	Token Token
	Tree  []TokenTree
}

type Valuer interface {
	getValue() interface{}
}

func (i Ident) getValue() interface{}  { return i }
func (n Number) getValue() interface{} { return n }
func (l LParen) getValue() interface{} { return l }
func (r RParen) getValue() interface{} { return r }

func main() {
	input := "Hello world"
	str, sourced := source(input)
	fmt.Printf("source: %#v, parsed: %#v, %#v\n", input, str, sourced)

	input = "(123 456) world"
	str, sourced = source(input)
	fmt.Printf("source: %#v, parsed: %#v, %#v\n", input, str, sourced)

	input = "((car cdr) cdr)"
	str, sourced = source(input)
	fmt.Printf("source: %#v, parsed: %#v, %#v\n", input, str, sourced)

	input = "()())))((()))"
	str, sourced = source(input)
	fmt.Printf("source: %#v, parsed: %#v, %#v\n", input, str, sourced)
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

func source(input string) (string, TokenTree) {
	var tokens []TokenTree
	for len(input) > 0 {
		nextInput, token := token(input)
		input = nextInput
		switch token.(type) {
		case LParen:
			var tt TokenTree
			input, tt = source(input)
			tokens = append(tokens, tt)
		case RParen:
			return input, TokenTree{Tree: tokens}
		default:
			tokens = append(tokens, TokenTree{Token: token})
		}
	}
	return input, TokenTree{Tree: tokens}
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
