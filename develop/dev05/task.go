package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

//fileBytes тип представляет собой содержимое текстового файла
type fileBytes [][]byte

//readFileBytes считывает данные из текстовых файлов.
//Принимает набор строк - путей к файлам, возвращает fileBytes и ошибку считывания
func readFilesBytes(filepaths ...string) (fileBytes, error) {

	lines := make([][]byte, 0)

	for _, path := range filepaths {
		content, err := ioutil.ReadFile(path)

		if err != nil {
			return nil, err
		}

		lines = append(lines, bytes.Split(content, []byte{'\n'})...)

	}
	res := fileBytes(lines)
	return res, nil
}

//getMatchResultIndexes возвращает индексы строк (слайсов байт) в fileBytes,
//удовлетворяющих паттерну (первый слайс), неудовлятворяющих паттерну (второй слайс)
//и ошибку матчинга в regexp
func getMatchResultIndexes(fb fileBytes, pattern string) (match []int, nomatch []int, err error) {
	for i := 0; i < len(fb); i++ {
		var ok bool
		ok, err = regexp.Match(pattern, fb[i])

		if err != nil {
			match, nomatch = nil, nil
			return
		}
		if ok {
			match = append(match, i)
		} else {
			nomatch = append(nomatch, i)
		}
	}

	return
}

//addContext добавляет к каждому индексу дополнительное количество
//before индексов до него и after индексов после него.
//Принимает: слайс индексов для которых будет добавляться контекст,
//длину (количество элементов) fileBytes, количество before индексов,
//количество after индексов. Возвращает: слайс индексов с контекстными индексами.
func addContext(in []int, lenfb int, before, after int) []int {
	cap := len(in) + len(in)*(before+after)
	result := make([]int, 0, cap)
	m := make(map[int]bool)

	for _, el := range in {
		start := el - before
		if start < 0 {
			start = 0
		}
		end := el + after
		if end >= lenfb {
			end = lenfb - 1
		}
		for j := start; j <= end; j++ {
			if !m[j] {
				result = append(result, j)
				m[j] = true
			}
		}

	}
	return result
}

//printFB распечатывает строки fileBytes, индексы которых совпадают
//с индексами входного слайса.
//Принимает: fileBytes, слайс индексов, флаг numOnly - требуется ли
//печатать только индексы строк, а не их значения.
func printFB(fb fileBytes, indxs []int, numOnly bool) {
	//создадим вспомогательную мапу для дальнейшей проверки вхождения индекса
	m := make(map[int]bool)
	for i := range indxs {
		m[indxs[i]] = true
	}

	for i := 0; i < len(fb); i++ {
		if m[i] {
			if numOnly {
				fmt.Println(i + 1)
			} else {
				fmt.Println(string(fb[i]))
			}
		}
	}
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
	lineNum := flag.Bool("n", false, "Print lines indexes")

	flag.Parse()

	//Проверяем наличие паттерна и имени хотя бы одного файла, как аргумента
	if flag.NArg() < 2 {
		log.Fatalf("File name and pattern required")
	}

	args := flag.Args()
	pattern := args[0]

	fileNames := args[1:]

	fb, err := readFilesBytes(fileNames...)

	if err != nil {
		log.Fatal(err.Error())
	}

	//Если требуется игнорировать case, изменим регулярное выражение
	if *ignore {
		pattern = "(?i)" + pattern
	}

	//Если требуется проверять совпадение строки целиков, изменим регулярное выражение
	if *fixed {
		pattern = "^" + pattern + "$"
	}

	matchIndx, noMatchIndx, err := getMatchResultIndexes(fb, pattern)

	if err != nil {
		log.Fatal(err)
	}

	//Context определяет after и before - колизию c разными значениями
	//устраняем, строго заполняя after и before значением из context
	if *context != 0 {
		*after = *context
		*before = *context
	}

	var printingIndxs []int

	if *invert {
		printingIndxs = addContext(noMatchIndx, len(fb), int(*before), int(*after))
	} else {
		printingIndxs = addContext(matchIndx, len(fb), int(*before), int(*after))
	}

	//Если требуется вывести только количество строк
	if *count {
		c := len(printingIndxs)
		fmt.Println(c)
		return
	}

	printFB(fb, printingIndxs, *lineNum)
}
