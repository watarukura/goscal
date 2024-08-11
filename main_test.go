package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestWhiteSpace(t *testing.T) {
	t.Run("white space", func(t *testing.T) {
		assert.Equal(t, "", whitespace("   "))
	})
}
func TestNumber(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		i, num := number("123.45 ")
		assert.Equal(t, " ", i)
		assert.Equal(t, Number("123.45"), num)
	})
}
func TestIdent(t *testing.T) {
	t.Run("ident", func(t *testing.T) {
		i, ident := ident("Adam")
		assert.Equal(t, "", i)
		assert.Equal(t, Ident("Adam"), ident)
	})
}
