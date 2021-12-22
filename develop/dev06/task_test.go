package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldListSet(t *testing.T) {
	input := "1,2,3"
	expected := fieldList([]int{1, 2, 3})

	var result fieldList
	err := result.Set(input)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, result))
}

func TestDelimStrings(t *testing.T) {
	input := []string{
		"aaa bb c",
		"aaabb c",
		"a aab bc",
	}
	expected := [][]string{
		{"aaa", "bb", "c"},
		{"aaabb", "c"},
		{"a", "aab", "bc"},
	}

	result := delimStrings(input, " ")

	assert.True(t, reflect.DeepEqual(result, expected))
}

func TestGetOnlyDelimeted(t *testing.T) {
	input := [][]string{
		{"aaa", "bb", "c"},
		{"aaabb", "c"},
		{"ccccac"},
		{"a", "aab", "bc"},
		{"bbababa"},
	}
	expected := [][]string{
		{"aaa", "bb", "c"},
		{"aaabb", "c"},
		{"a", "aab", "bc"},
	}

	result := getOnlyDelimeted(input)

	assert.True(t, reflect.DeepEqual(expected, result))
}
