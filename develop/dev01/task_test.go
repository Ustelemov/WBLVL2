package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//TestSystemNotEmpty check System time not equality ""
func TestSystemNotEmpty(t *testing.T) {
	timeString := getSystemUnixTime()
	assert.NotZero(t, timeString)
}

//TestSystemNotZero check System time not equality "0001-01-01 00:00:00 +0000 UTC"
func TestSystemTimeNotZero(t *testing.T) {
	timeString := getSystemUnixTime()
	assert.NotEqual(t, timeString, time.Time{}.Format(time.UnixDate))
}

//TestNetTimeNotZero check Net time not equality ""
func TestNetTimeNotZero(t *testing.T) {
	server := "time.apple.com"
	timeString, _ := getNetUnixTime(server)
	assert.NotZero(t, timeString)
}

//TestSystemNotZero check Net time not equality "0001-01-01 00:00:00 +0000 UTC"
func TestNetTimeNotEmpty(t *testing.T) {
	server := "time.apple.com"
	timeString, _ := getNetUnixTime(server)
	assert.NotEqual(t, timeString, time.Time{}.Format(time.UnixDate))
}

//TestNetTimeNotError check Net time err being nil
func TestNetTimeNotError(t *testing.T) {
	server := "time.apple.com"
	_, err := getNetUnixTime(server)
	assert.Nil(t, err)
}
