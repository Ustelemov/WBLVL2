package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

/*
=== Утилита grep ===
Реализовать утилиту фильтрации (man grep)
Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//isMatch проверяет соответствие слайса байт регулярному выражению (паттерну).
//Принимает слайс байт и паттерн.
//Возвращает результат соответствия и ошибку матчинга.
func isMatch(b []byte, pattern string) (bool, error) {
	return regexp.Match(pattern, b)
}

//appendBuffer добавляет в буфер последних элементов новый элемент
//(если буфер уже заполнен добавляемый элемент "выбивает" самый первый элемент)
//Принимает ссылку на слайс строк-cтруктур (модифицирует его в ходе выполнения)
//и строку, которую следует добавить.
func appendBuffer(buffer *[]line, el line) {
	if cap(*buffer) == 0 {
		return
	}
	if len(*buffer) < cap(*buffer) {
		*buffer = append(*buffer, el)
	} else {
		*buffer = append((*buffer)[1:len(*buffer)], el)
	}
}

//структура line определяет текстовую строку в файле\консоле
//содержит саму строку и её порядковый номер
type line struct {
	text string
	num  int
}

func main() {
	//Устанавливаем возможные флаги
	after := flag.Uint("A", 0, "Printing N lines after match")
	before := flag.Uint("B", 0, "Printing N lines before match")
	context := flag.Uint("C", 0, "Printing N lines before and after match")
	count := flag.Bool("c", false, "Printing count of matching lines")
	ignore := flag.Bool("i", false, "Ignore case of matching pattern")
	invert := flag.Bool("v", false, "Exclude string if matching")
	fixed := flag.Bool("F", false, "Exact equality to string not pattern")
	numering := flag.Bool("n", false, "Print lines indexes")

	flag.Parse()

	//Проверяем наличие паттерна и имени хотя бы одного файла (или -), как аргумента
	if flag.NArg() < 2 {
		log.Fatalf("File name or - for console input and pattern required")
	}

	args := flag.Args()
	pattern := args[0]

	fileName := args[1]

	//Context определяет after и before - колизию c разными значениями
	//устраняем, строго заполняя after и before значением из context
	if *context != 0 {
		*after = *context
		*before = *context
	}

	//Если требуется игнорировать case, изменим регулярное выражение
	if *ignore {
		pattern = "(?i)" + pattern
	}

	//Если требуется проверять совпадение строки целиков, изменим регулярное выражение
	if *fixed {
		pattern = "^" + pattern + "$"
	}

	isConsole := fileName == "-"
	var scanner *bufio.Scanner

	//Зададим источник данных: консоль или файл
	if isConsole {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	//создадим буфер с емкостью before
	bufferBefore := make([]line, 0, *before)
	//строк осталось вывести после данной
	afterCount := 0
	//номер строки (начнем с 1)
	lineNumber := 0
	//количество совпадений
	matches := 0

	//map, сохраняющий номера напечатанных строки с true
	printedMap := make(map[line]bool, 0)

	for {
		lineNumber++
		flag := scanner.Scan()
		text := scanner.Text()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		//EOF
		if !flag {
			break
		}

		ok, err := isMatch([]byte(text), pattern)
		if err != nil {
			log.Fatal(err)
		}

		//XOR - изменим значение наоборот, если требуется инвертировать
		ok = ok != *invert
		//строки, которые потребуется напечатать в данные проход
		//сюда попадут ещё before строки
		linesToPrint := make([]line, 0)

		if ok {
			matches++
			//добавим нужные элементы перед
			linesToPrint = append(linesToPrint, bufferBefore...)
			//обновим, так как нашли матчинг
			bufferBefore = make([]line, 0, *before)
			afterCount = int(*after)

			linesToPrint = append(linesToPrint, line{text: text, num: lineNumber})
		} else {
			//если для последнего матчинга напечатали не все строки после,
			//то печатаем и уменьшаем счетчик
			if afterCount != 0 {
				afterCount--
				linesToPrint = append(linesToPrint, line{text: text, num: lineNumber})
			}
			//запишим строку и before буфер
			appendBuffer(&bufferBefore, line{text: text, num: lineNumber})
		}

		if isConsole && *count { //в консоле - результат пустой
			fmt.Println("")
		} else if !*count { //если не файл с -c флагом
			for _, v := range linesToPrint {
				if !printedMap[v] { //если строку ещё не печатали
					printedMap[v] = true
					if *numering {
						fmt.Fprintf(os.Stdout, "%d:%s\n", v.num, v.text)
					} else {
						fmt.Fprintln(os.Stdout, text)
					}
				}
			}
		}
	}
	if !isConsole && *count {
		fmt.Fprintln(os.Stdout, matches)
	}
}
