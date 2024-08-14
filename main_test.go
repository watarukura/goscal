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
func TestExpression(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want Expression
	}{
		{name: "valid ident", arg: "Adam", want: Ident("Adam")},
		{name: "valid num", arg: "123", want: Number(123)},
		{name: "valid add", arg: "123.4 + 234.5", want: Add{Number(123.4), Number(234.5)}},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, actual, _ := expr(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
