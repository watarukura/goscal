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
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := Parse(test.arg)
			assert.Equal(t, test.want, actual)
		})
	}
}
