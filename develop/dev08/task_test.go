package main

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCD(t *testing.T) {
	cd := cdCMD{}
	_, err := cd.exec([]string{".."}, nil, false)

	assert.Nil(t, err)
}

func TestEcho(t *testing.T) {
	echo := echoCMD{}
	var expected bytes.Buffer
	expected.WriteString("hello\n")

	result, err := echo.exec([]string{"hello"}, nil, false)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, *result))
}

func TestFork(t *testing.T) {
	fork := forkCMD{}
	var expected bytes.Buffer
	expected.WriteString("hello\n")

	result, err := fork.exec([]string{"echo", "hello"}, nil, false)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, *result))
}

func TestPWD(t *testing.T) {
	pwd := pwdCMD{}

	result, err := pwd.exec(nil, nil, false)

	assert.Nil(t, err)
	assert.True(t, len(result.String()) > 0)
}

func TestPS(t *testing.T) {
	ps := psCMD{}

	result, err := ps.exec(nil, nil, false)

	assert.Nil(t, err)
	assert.True(t, len(result.String()) > 0)
}

func TestExecCommand(t *testing.T) {
	var expected bytes.Buffer
	expected.WriteString("hello\n")

	result, err := execCommand("echo hello", nil, false)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expected, *result))
}

func TestExecCommands(t *testing.T) {
	var expected bytes.Buffer
	expected.WriteString("hello\n")

	err := execCommands("echo hello | grep h")

	assert.Nil(t, err)
}
