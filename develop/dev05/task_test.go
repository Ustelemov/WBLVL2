package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMatchOk(t *testing.T) {
	input := []byte("qweasd")
	pattern := "qwe"

	ok, err := isMatch(input, pattern)
	assert.Nil(t, err)
	assert.True(t, ok)

}

func TestIsMatchNotOK(t *testing.T) {
	input := []byte("qweasd")
	pattern := "qweee"

	ok, err := isMatch(input, pattern)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestAppendBuffer(t *testing.T) {
	buf := make([]line, 0, 2)
	expected := []line{{"3", 3}, {"4", 4}}
	appendBuffer(&buf, line{"1", 1})
	appendBuffer(&buf, line{"2", 2})
	appendBuffer(&buf, line{"3", 3})
	appendBuffer(&buf, line{"4", 4})

	assert.True(t, reflect.DeepEqual(expected, buf))
}
