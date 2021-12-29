Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
https://www.pixelstech.net/article/1542724558-Behavior-of-defer-function-in-named-return-function
Вывод: 2 1
Выражение defer добавляет вызов функции, для которой оно объявлено, в стек вызовов.
Извлечение вызовов из стека происходит в порядке LIFO (Last in First Out).
defer функции вызываются непосредственно до возврата из ближайшей окружающей
функции.

Разница двух функций здесь в том, что в случае с неименованным возвратом
при вызове выражения return происходит копирование значения x в область
памяти из которой будет произведен возврат значения при возврата из функции
(и только после этого вызывается defer).

В случае же именованного возврата при вызове return копирования не происходит, блок defer изменяет значение переменной и при возврате из
функции будет возвращено значение из области памяти возвращаемой именнованой переменной.
```