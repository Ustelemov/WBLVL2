package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLower(t *testing.T) {
	input := []string{"Aa", "A", "AB", "ab", "a", "C", "AA"}
	expected := []string{"aa", "a", "ab", "ab", "a", "c", "aa"}

	result := toLower(input)
	assert.True(t, reflect.DeepEqual(result, expected))
}

func TestDeleteRepeated(t *testing.T) {
	input := []string{"aa", "a", "ab", "ab", "a", "c", "aa"}
	expected := []string{"aa", "a", "ab", "c"}

	result := deleteRepeated(input)
	assert.True(t, reflect.DeepEqual(result, expected))
}

func TestMakeAnagrammDict(t *testing.T) {
	input := []string{"орел", "Катер", "Актер", "рысь", "сырь", "катер", "Сырь", "рысь", "катер", "терка"}
	expected := map[string][]string{
		"катер": {"актер", "катер", "терка"},
		"рысь":  {"рысь", "сырь"},
	}

	result := makeAnagrammDict(input)
	assert.True(t, reflect.DeepEqual(expected, result))
}
