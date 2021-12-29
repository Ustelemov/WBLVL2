Что выведет программа? Объяснить вывод программы.

```go
package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
```

Ответ:
```
(https://go.dev/tour/moretypes/7)

Результатом будет вывод: [77 78 79].
Здесь в переменную b будет записан слайс, полученный "слайсингом" [1:4] массива, где 1 - это индекс первого включаемого элемента, а 4 - индекс обозначющий правую невключаемую границу.
То есть элементы с индексами 1, 2, 3 исходного массива попадут в результат.
b будет представлять собой слайс со ссылкой на массив а, длинной 3 и емкостью 4.
```