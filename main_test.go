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
		assert.Equal(t, " ", number("123.45 "))
	})
}
func TestIdent(t *testing.T) {
	t.Run("ident", func(t *testing.T) {
		assert.Equal(t, "", ident("Adam"))
	})
}
