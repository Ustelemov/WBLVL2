package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===
Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)
В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.
Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.

Синтаксический разбор упакованной строки:
String: []RuneRepeat
RuneRepeat: Element Number
Element: Slash Rune | Letter
Slash
Number
Rune
Letter

Реализуемые ниже функции соответствуют данному синтаксическому разбору
*/

//RuneRepeat структура, содержащая руну и количество её повторений
type runeRepeat struct {
	value rune
	times int
}

//tryUnpackString принимает входную запакованную строку и возвращает
//строку результата и bool значение успешности распаковки
//для проверки коррекности распаковки использовать comma ok
func tryUnpackString(input string) (string, bool) {
	var sb strings.Builder
	runes := []rune(input) //переводим строку в слайс рун

	//Функции нижних уровней будут уменьшать слайс, цикл пока он не пустой
	for len(runes) > 0 {
		runeRepeat, res := tryParseRuneRepeat(&runes)
		if !res {
			return "", false
		}

		//Повторяем элемент нужное количество раз и сохраняем в билдер
		sb.WriteString(strings.Repeat(string(runeRepeat.value), runeRepeat.times))
	}

	return sb.String(), true
}

//tryParseRuneRepeat пытается распрасить элемент RuneRepeat
//принимает ссылку на слайс рун
//возвращает ссылку на структуру RuneRepeate и результат парсинга bool
func tryParseRuneRepeat(runes *[]rune) (*runeRepeat, bool) {
	if elem, res := tryParseElement(runes); res {
		if number, res := tryParseNumber(runes); res {
			return &runeRepeat{
				value: elem,
				times: number,
			}, true
		}
		return &runeRepeat{
			value: elem,
			times: 1, //если числа повторений нет, значит оно = 1
		}, true
	}

	return nil, false
}

//tryParseElement пытается распарсить Element
//принимает ссылку на слайс рун
//и возвращает rune и результат парсинга bool
func tryParseElement(runes *[]rune) (rune, bool) {
	if _, res := tryParseSlash(runes); res {
		if elem, res := tryParseRune(runes); res {
			return elem, true
		}
	} else if elem, res := tryParseLetter(runes); res {
		return elem, true
	}

	return rune(-1), false

}

//tryParseNumber пытается распарсить Number
//принимает ссылку на слайс рун (модифицирует ссылку, удаляя распаршенный
//элемент в случае успеха парсинга)
//возвращает целое число int и результат парсинга bool
func tryParseNumber(runes *[]rune) (int, bool) {
	//Если слайс пустой - парсить нечего
	if len(*runes) == 0 {
		return -1, false
	}

	if numRunes, err := strconv.Atoi(string((*runes)[0])); err == nil {
		*runes = (*runes)[1:]
		return numRunes, true
	}

	return -1, false

}

//tryParseSlash пытается распарсить Slash
//принимает ссылку на слайс рун (модифицирует ссылку, удаляя распаршенный
//элемент в случае успеха парсинга)
//возвращает rune (просто заглушку) и результат парсинга bool
func tryParseSlash(runes *[]rune) (rune, bool) {
	//Если слайс пустой - парсить нечего
	if len(*runes) == 0 {
		return rune(-1), false
	}

	if (*runes)[0] == '\\' {
		*runes = (*runes)[1:]
		return '\\', true
	}

	return rune(-1), false
}

//tryParseRune пытается распарсить Rune
//принимает ссылку на слайс рун (модифицирует ссылку, удаляя распаршенный
//элемент в случае успеха парсинга)
//возвращает rune и результат парсинга bool
func tryParseRune(runes *[]rune) (rune, bool) {
	//Если слайс пустой - парсить нечего
	if len(*runes) == 0 {
		return rune(-1), false
	}

	result := (*runes)[0]
	*runes = (*runes)[1:]

	return result, true
}

//tryParseLetter пытается распарсить Letter
//принимает ссылку на слайс рун (модифицирует ссылку, удаляя распаршенный
//элемент в случае успеха парсинга)
//возвращает rune и результат парсинга bool
func tryParseLetter(runes *[]rune) (rune, bool) {
	//Если слайс пустой - парсить нечего
	if len(*runes) == 0 {
		return rune(-1), false
	}

	result := (*runes)[0]

	if unicode.IsLetter(result) {
		*runes = (*runes)[1:]
		return result, true
	}

	return rune(-1), false
}

func main() {
	//строку для корректной передаем "как есть", используя ``
	input := `q4w5\43`
	res, ok := tryUnpackString(input)

	fmt.Printf("Input string: %s\n", input)
	if !ok {
		fmt.Println("Unpacking error. Bad input")
	}
	fmt.Printf("Unpacking result: %s\n", res)
}
