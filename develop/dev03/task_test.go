package main

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filepath = "test.txt"

func TestReadFile(t *testing.T) {
	_, err := readFileBytes(filepath)
	assert.Nil(t, err)
}

func TestWriteFile(t *testing.T) {
	readed, err := readFileBytes(filepath)
	assert.Nil(t, err)

	err = writeFileBytes(filepath, readed)
	assert.Nil(t, err)

	readedNew, err := readFileBytes(filepath)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(readed, readedNew))
}

func TestSortBytes(t *testing.T) {
	input := [][]byte{
		{8, 2, 3},
		{4, 5, 6},
		{1, 2, 1},
	}

	expected := [][]byte{
		{1, 2, 1},
		{4, 5, 6},
		{8, 2, 3},
	}

	fbin := fileBytes(input)
	fbexp := fileBytes(expected)

	getF := buildGetF()
	compF := buildCompF()

	sortBytes(fbin, getF, compF)

	assert.True(t, reflect.DeepEqual(fbin, fbexp))
}

func TestBuildGetFColumn(t *testing.T) {
	inputS := "a b c\nq w e"
	fb := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))
	column := 2

	getF := buildGetFColumn(column)
	b1, b2 := getF(fb, 0, 1)
	assert.True(t, reflect.DeepEqual(b1, []byte{'c'}))
	assert.True(t, reflect.DeepEqual(b2, []byte{'e'}))
}

func TestBuildGetF(t *testing.T) {
	inputS := "a b c\nq w e"
	fb := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))

	getF := buildGetF()
	b1, b2 := getF(fb, 0, 1)
	assert.True(t, reflect.DeepEqual(b1, []byte("a b c")))
	assert.True(t, reflect.DeepEqual(b2, []byte("q w e")))
}

func TestBuildCompF(t *testing.T) {
	inputS := "a b c\nq w e"
	fb := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))

	getF := buildGetF()
	b1, b2 := getF(fb, 0, 1)

	compF := buildCompF()
	b1Lessb2 := compF(b1, b2)

	assert.True(t, b1Lessb2)
}

func TestBuildCompFloat(t *testing.T) {
	inputS := "1.77\n1.33"
	fb := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))

	getF := buildGetF()
	b1, b2 := getF(fb, 0, 1)

	compF := buildCompFloat()
	b1Lessb2 := compF(b1, b2)

	assert.False(t, b1Lessb2)
}

func TestReverse(t *testing.T) {
	inputS := "1.77\n1.33"
	fbIn := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))

	expectedS := "1.33\n1.77"
	fbExp := fileBytes(bytes.Split([]byte(expectedS), []byte{'\n'}))

	reverse(fbIn)

	assert.True(t, reflect.DeepEqual(fbExp, fbIn))
}

func TestRemoveRepeats(t *testing.T) {
	inputS := "1.77\n1.33\n1.77\n1.33"
	fbIn := fileBytes(bytes.Split([]byte(inputS), []byte{'\n'}))

	expectedS := "1.77\n1.33"
	fbExp := fileBytes(bytes.Split([]byte(expectedS), []byte{'\n'}))

	fbRes := removeRepeats(fbIn)

	assert.True(t, reflect.DeepEqual(fbExp, fbRes))
}
