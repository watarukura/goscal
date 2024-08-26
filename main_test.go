package main

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestParse(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want float64
	}{
		{name: "valid add", arg: "123.4 + 234.5", want: 357.9},
		{name: "valid paren", arg: "(123.4 + 234.5)", want: 357.9},
		{name: "valid add 3 term", arg: "1 + 2 + 3", want: 6},
		{name: "valid nested paren", arg: "((1 + 2) + (3 + 4)) + 5 + 6", want: 21},
		{name: "valid unary op", arg: "-100", want: -100},
		{name: "valid pi", arg: "pi", want: 3.141592653589793},
		{name: "valid add pi", arg: "(123 + 456 ) + pi", want: 582.1415926535898},
		{name: "valid mul", arg: "100 * -10", want: -1000},
		{name: "valid paren div", arg: "(100 + 10) / 10", want: 11},
		{name: "valid sqrt", arg: "sqrt(100)", want: 10},
		{name: "valid pow", arg: "pow(3, 2)", want: 9},
		{name: "valid pow + pow", arg: "pow(3, 2) + pow(3, 2)", want: 18},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := Parse(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
