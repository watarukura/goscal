package main

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestWhiteSpace(t *testing.T) {
	t.Run("white space", func(t *testing.T) {
		assert.Equal(t, "", whitespace("   "))
	})
}
func TestNumber(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Expression
	}{
		{name: "valid", arg: "123.45 ", want: Number(123.45)},
		{name: "not num", arg: "aaa", want: Expression(nil)},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := number(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
func TestIdent(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Expression
	}{
		{name: "valid", arg: "Adam ", want: Ident("Adam")},
		{name: "not identifier", arg: "123", want: Expression(nil)},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := ident(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
func TestAdd(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Expression
	}{
		{name: "valid add", arg: "123.4 + 234.5", want: Add{Number(123.4), Number(234.5)}},
		{name: "valid add 3 term", arg: "1 + 2 + 3", want: Add{Add{Number(1), Number(2)}, Number(3)}},
		{name: "valid paren", arg: "(1 + 2) + 3", want: Add{Add{Number(1), Number(2)}, Number(3)}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := add(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
func TestExpr(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Expression
	}{
		{name: "valid add", arg: "123.4 + 234.5", want: Add{Number(123.4), Number(234.5)}},
		{name: "valid paren", arg: "(123.4 + 234.5)", want: Add{Number(123.4), Number(234.5)}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := expr(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
func TestEval(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Number
	}{
		{name: "valid add", arg: "123.4 + 234.5", want: Number(357.9)},
		{name: "valid paren", arg: "(123.4 + 234.5)", want: Number(357.9)},
		{name: "valid pi", arg: "pi", want: Number(3.141592653589793)},
		{name: "valid add 3 term", arg: "1 + 2 + 3", want: Number(6)},
		{name: "valid nested paren", arg: "((1 + 2) + (3 + 4)) + 5 + 6", want: Number(21)},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, expr, _ := expr(test.arg)
			actual := eval(expr)
			assert.Equal(t, test.want, actual)
		})
	}
}
