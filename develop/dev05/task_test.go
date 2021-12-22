package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filePath = "test.txt"

func TestReadFilesBytes(t *testing.T) {
	_, err := readFilesBytes(filePath)
	assert.Nil(t, err)
}

func TestMatchResultIndexes(t *testing.T) {
	input := [][]byte{
		[]byte("qweasd"),
		[]byte("zxcasdqwwe"),
		[]byte("qwsfsde"),
		[]byte("123rqwe"),
	}

	expectedMatch := []int{0, 3}
	expectedNoMatch := []int{1, 2}

	pattern := "qwe"

	match, nomatch, err := getMatchResultIndexes(fileBytes(input), pattern)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedMatch, match))
	assert.True(t, reflect.DeepEqual(expectedNoMatch, nomatch))
}

func TestAddContext(t *testing.T) {
	input := []int{3, 9}
	lenfb := 10
	before, after := 2, 2

	expect := []int{1, 2, 3, 4, 5, 7, 8, 9}

	result := addContext(input, lenfb, before, after)

	assert.True(t, reflect.DeepEqual(expect, result))

}
