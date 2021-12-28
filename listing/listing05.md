Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Результат будет: error.
Переменная типа интерфейс error будет хранить (*main.customError)(nil),
где тип значения отличен от nil. Поэтому как и в третьем задании
сравнение с nil будет выдавать false.
```