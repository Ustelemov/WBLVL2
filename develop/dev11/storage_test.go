package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventStore(t *testing.T) {
	res := NewEventStore()
	assert.NotNil(t, res)
}

func TestSave(t *testing.T) {
	events := NewEventStore()

	events.Save("123", time.Now())
	events.Save("123", time.Now())

	assert.Equal(t, len(events.m), 2)
}

func TestLoad(t *testing.T) {
	events := NewEventStore()
	events.Save("123", time.Now())

	event, ok := events.Load(1)
	assert.True(t, ok)
	assert.Equal(t, event.ID, 1)
}

func TestChange(t *testing.T) {
	events := NewEventStore()
	events.Save("123", time.Now())

	ok, err := events.Change(1, "1234", time.Time{})
	assert.True(t, ok)
	assert.Nil(t, err)

	event, ok := events.Load(1)
	assert.True(t, ok)
	assert.Equal(t, event.Text, "1234")
}

func TestGetTodays(t *testing.T) {
	events := NewEventStore()
	events.Save("123", time.Now())
	events.Save("123", time.Now())
	events.Save("123", time.Date(2021, 12, 12, 0, 0, 0, 0, time.Now().Location()))
	events.Save("123", time.Date(2021, 27, 12, 0, 0, 0, 0, time.Now().Location()))

	result, err := events.GetTodays()
	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
}

func TestGetThisWeeks(t *testing.T) {
	events := NewEventStore()
	events.Save("123", time.Now())
	events.Save("123", time.Now())
	events.Save("123", time.Date(2021, 05, 12, 0, 0, 0, 0, time.Now().Location()))
	events.Save("123", time.Date(2021, 01, 12, 0, 0, 0, 0, time.Now().Location()))

	result, err := events.GetThisWeeks()
	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
}

func TestGetThisMonths(t *testing.T) {
	events := NewEventStore()
	events.Save("123", time.Now())
	events.Save("123", time.Now())
	events.Save("123", time.Date(2021, 05, 11, 0, 0, 0, 0, time.Now().Location()))
	events.Save("123", time.Date(2021, 01, 11, 0, 0, 0, 0, time.Now().Location()))

	result, err := events.GetThisMonths()
	assert.Nil(t, err)
	assert.Equal(t, len(result), 2)
}

func TestGetBetween(t *testing.T) {
	events := NewEventStore()
	start := time.Now()
	between := time.Now()
	end := time.Now()

	events.Save("123", between)
	events.Save("123", time.Date(2021, 05, 11, 0, 0, 0, 0, time.Now().Location()))
	events.Save("123", time.Date(2021, 01, 11, 0, 0, 0, 0, time.Now().Location()))

	result, err := events.getBetween(start, end)
	assert.Nil(t, err)
	assert.Equal(t, len(result), 1)
}

func TestInTimeSpan(t *testing.T) {
	start := time.Now()
	between := time.Now()
	end := time.Now()

	result := inTimeSpan(start, end, between)
	assert.True(t, result)
}
