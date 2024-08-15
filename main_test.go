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
		{name: "valid hello world", arg: "Hello + World", want: Add{Ident("Hello"), Ident("World")}},
		{name: "valid paren world", arg: "(123 + 456 ) + World", want: Add{Add{Number(123), Number(456)}, Ident("World")}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := expr(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
