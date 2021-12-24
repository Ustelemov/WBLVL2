package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {

	const minSec = 5
	start := time.Now()

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(minSec*time.Second), //минимальный по времени done-канал
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	durSec := time.Since(start).Seconds()

	assert.True(t, durSec < minSec*1.05)
}
