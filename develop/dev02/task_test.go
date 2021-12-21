package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLetterOK(t *testing.T) {
	input := []rune("a")
	_, ok := tryParseLetter(&input)
	assert.True(t, ok)
}

func TestParseLetterBad(t *testing.T) {
	input := []rune("1")
	_, ok := tryParseLetter(&input)
	assert.False(t, ok)
}

func TestParseNumberOK(t *testing.T) {
	input := []rune("1")
	_, ok := tryParseNumber(&input)
	assert.True(t, ok)
}

func TestParseNumberBad(t *testing.T) {
	input := []rune("a")
	_, ok := tryParseNumber(&input)
	assert.False(t, ok)
}

func TestParseSlashOK(t *testing.T) {
	input := []rune(`\`)
	_, ok := tryParseSlash(&input)
	assert.True(t, ok)
}

func TestParseSlashBad(t *testing.T) {
	input := []rune("a")
	_, ok := tryParseSlash(&input)
	assert.False(t, ok)
}

func TestParseRuneOK(t *testing.T) {
	input := []rune("a")
	_, ok := tryParseRune(&input)
	assert.True(t, ok)
}

func TestUnpackEmptyOK(t *testing.T) {
	input := ""
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, "")
}

func TestUnpackSimpleOK(t *testing.T) {
	input := "a4bc2d5e"
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, "aaaabccddddde")
}

func TestUnpackSimple1OK(t *testing.T) {
	input := "abcd"
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, "abcd")
}

func TestUnpackSimpleBad(t *testing.T) {
	input := "45"
	_, ok := tryUnpackString(input)
	assert.False(t, ok)
}

func TestUnpackEscapedOK(t *testing.T) {
	input := `qwe\4\5`
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, "qwe45")
}

func TestUnpackEscaped1OK(t *testing.T) {
	input := `qwe\45`
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, "qwe44444")
}

func TestUnpackEscaped2OK(t *testing.T) {
	input := `qwe\\5`
	res, ok := tryUnpackString(input)
	assert.True(t, ok)
	assert.Equal(t, res, `qwe\\\\\`)
}
