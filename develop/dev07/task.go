package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===
Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.
Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}
Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}
start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)
fmt.Printf(“fone after %v”, time.Since(start))
*/

//or объединяет набор done каналов только для чтения в
//единственный done канал. Результирующий канал
//будет закрыт, если один из каналов входного набора будет закрыт.
func or(channels ...<-chan interface{}) <-chan interface{} {
	res := make(chan interface{})

	go func() {
		for {
			for _, c := range channels {
				select {
				case _, ok := <-c:
					if !ok { //если канал закрыт
						close(res) //закрываем single канал
						return     //и выходим из анонимной функции
					}
				default: //если канал не закрыт, проверяем следующий
					continue
				}
			}
		}
	}()

	return res
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(5*time.Second), //минимальный по времени done-канал
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
