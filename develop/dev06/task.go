package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===
Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные
Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//Тип, используемый для парсинга флага fields, реализующий Value interface
type fieldList []int

//String - реализация интерфейса Value.
//Представление объекта в виде строки.
//Возвращает - строку.
func (f *fieldList) String() string {
	return fmt.Sprintln(*f)
}

//Set - реализация интерфейса Value.
//Принимает строку и устанавливает значения объекта из строки:
//заполняет слайс индексов, разделенных в исходной строке запятой.
//Возвращает error.
func (f *fieldList) Set(s string) error {
	splited := strings.Split(s, ",")
	for _, v := range splited {
		el, err := strconv.Atoi(v)

		if err != nil {
			return err
		}

		*f = append(*f, el)
	}
	return nil
}

//delimStrings принимает слайс строк и разделитель.
//Выдает: слайс слайсов строк, где каждый внутренний слайс - строки,
//разбитые разделителем.
func delimStrings(in []string, d string) [][]string {
	var result [][]string
	for _, s := range in {
		result = append(result, strings.Split(s, d))
	}
	return result
}

//getDelimeted возвращает строки, которые были разделены разделителем.
//Принимает - слайс слайсов строк.
//Возвращает - слайс слайсов строк, где внутренние слайсы с количеством
//элементов больше 1 (так как были разделены разделителем)
func getOnlyDelimeted(in [][]string) [][]string {
	var result [][]string
	for _, v := range in {
		if len(v) != 1 {
			result = append(result, v)
		}
	}
	return result
}

//printStrings распечатывает строки с конкретными индексами cтобцов
//(столбцы образовались в результате деления строки с помощью делителя),
//где индексы задаются входным параметром слайсом индексов.
//Принимает: слайс слайсов строк, слайс индексов и разделитель для печати.
func printStrings(in [][]string, fl fieldList, delim string) {
	//Создадим дополнительную мапу для удобства проверки вхождения индекса
	m := make(map[int]bool, 0)
	for i := range fl {
		m[fl[i]] = true
	}

	for i := range in {
		f := true //первый столбик
		for j := range in[i] {
			if m[j] || len(fl) == 0 { //если столбцы не указаны - печатаем все
				if f {
					f = false
				} else {
					fmt.Print(delim) //разделитель перед непервым столбцов
				}
				fmt.Print(in[i][j])
			}
		}
		if !f { //если что-то было напечатано в строке
			fmt.Println("") //новая строка
		}
	}
}

func main() {

	var fields fieldList

	//Установим флаги
	flag.Var(&fields, "f", "Select columns number comma separated")
	delim := flag.String("d", "\t", "Delimeter for split string on columns")
	separated := flag.Bool("s", false, "Print lines only with delimeter")
	flag.Parse()

	//Если помимо флагов было добавлено что-то другое
	if flag.NArg() > 0 {
		log.Fatal("Bad input: no args should be provided")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		lines := *new([]string) // каждый цикл обновляем ссылку

		scanner.Scan()
		text := scanner.Text()

		lines = strings.Split(text, "\n")

		delimed := delimStrings(lines, *delim)
		if *separated { //если нужны только разделенные делителем строки
			delimed = getOnlyDelimeted(delimed)
		}

		printStrings(delimed, fields, *delim)
	}
}
