package main

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

//fileBytes тип представляет собой содержимое текстового файла
type fileBytes [][]byte

//readFileBytes считывает данные из текстового файла.
//Принимает строку пути к файлу, возвращает fileBytes и ошибку считывания
func readFileBytes(filepath string) (fileBytes, error) {
	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	lines := bytes.Split(content, []byte{'\n'})

	res := fileBytes(lines)

	return res, nil
}

//writeFileBytes записывает fileBytes в файл.
//Принимает строку пути к файлу и filebytes, возвращает ошибку записи
func writeFileBytes(filepath string, fb fileBytes) error {
	content := bytes.Join(fb, []byte{'\n'})
	return ioutil.WriteFile(filepath, content, 0644)
}

//sortBytes сортирует слайс.
//Принимает: fileBytes - исходные данные для сортировки, изменяются в ходе сортировки,
//getF - функцию получения двух слайсов байт, которые будут сравниваться
//compF - функция Less для двух слайсов байт
func sortBytes(in fileBytes,
	getF func(fileBytes, int, int) ([]byte, []byte),
	compF func([]byte, []byte) bool) {
	sort.Slice(in, func(i, j int) bool {
		return compF(getF(in, i, j))
	})
}

//buildGetFColumn создает функцию получения элементов строк для сравнения
//по столбцу. Принимает номер столбца, возвращает функцию получения элементов
func buildGetFColumn(column int) func(fileBytes, int, int) ([]byte, []byte) {
	return func(in fileBytes, i, j int) ([]byte, []byte) {
		wordsI := strings.Fields(string(in[i]))
		wordsJ := strings.Fields(string(in[j]))

		resI := make([]byte, 0)
		resJ := make([]byte, 0)

		if len(wordsI) > column {
			resI = []byte(wordsI[column])
		}

		if len(wordsJ) > column {
			resJ = []byte(wordsJ[column])
		}

		return resI, resJ
	}
}

//buildGetFColumn создает базовую функцию получения элементов строк для сравнения.
//Возвращает функцию получения элементов
func buildGetF() func(fileBytes, int, int) ([]byte, []byte) {
	return func(in fileBytes, i, j int) ([]byte, []byte) {
		return in[i], in[j]
	}
}

//buildCompF создает базовую функцию сравнения (Less) слайсов байт
//Возвращает функцию сравнения (Less)
func buildCompF() func([]byte, []byte) bool {
	return func(b1 []byte, b2 []byte) bool {
		return bytes.Compare(b1, b2) == -1
	}
}

//buildCompF создает базовую функцию сравнения (Less) слайсов байт
//на основе сравнения float значениях, содержащихся в слайсах.
//Возвращает функцию сравнения (Less)
//Если парсинг float значений неудачен - выход из программы с кодом 1
func buildCompFloat() func([]byte, []byte) bool {
	return func(b1 []byte, b2 []byte) bool {
		if len(b1) == 0 && len(b2) == 0 {
			return false
		}

		if len(b1) == 0 || len(b2) == 0 {
			return len(b1) == 0
		}

		floatB1, err := strconv.ParseFloat(string(b1), 64)

		if err != nil {
			log.Fatalf(err.Error())
		}

		floatB2, err := strconv.ParseFloat(string(b2), 64)

		if err != nil {
			log.Fatalf(err.Error())
		}

		return floatB1 < floatB2
	}
}

//reverse переворачивает строки файла в обратном порядке
//Принимает fileBytes, изменяя его
func reverse(in fileBytes) {
	for i, j := 0, len(in)-1; i < j; i, j = i+1, j-1 {
		in[i], in[j] = in[j], in[i]
	}
}

//removeRepeats удаляет повторяющеся строки в файле
//Принимает fileBytes, возвращает новый fileBytes
func removeRepeats(in fileBytes) fileBytes {
	result := make([][]byte, 0)
	m := make(map[string]bool, 0)

	for _, el := range in {
		if !m[string(el)] {
			result = append(result, el)
			m[string(el)] = true
		}
	}

	return fileBytes(result)
}

func main() {
	//Устанавливаем возможные флаги
	sortingColumn := flag.Int("k", -1, "Column for sorting (word from line)")
	isNumberSorting := flag.Bool("n", false, "Sorting by number value")
	isReverseSorting := flag.Bool("r", false, "Reverse sorting")
	isRemoveRepeating := flag.Bool("u", false, "Remove repeating strings")
	flag.Parse()

	//Проверяем наличие имени файла, как аргумента
	if flag.NArg() < 1 {
		log.Fatalf("File name required")
	}

	fileName := flag.Arg(0)

	readResult, err := readFileBytes(fileName)

	if err != nil {
		log.Fatal(err.Error())
	}

	var getF func(fileBytes, int, int) ([]byte, []byte)
	var compF func(b1, b2 []byte) bool

	//Проверяем требуется ли сортировка по конкретному столбцу
	if *sortingColumn != -1 {
		getF = buildGetFColumn(*sortingColumn)
	} else {
		getF = buildGetF()
	}

	//Проверяем требуется ли сортировка по числову значению
	if *isNumberSorting {
		compF = buildCompFloat()
		if *sortingColumn == -1 {
			getF = buildGetFColumn(0)
		}

	} else {
		compF = buildCompF()
	}

	sortBytes(readResult, getF, compF)

	//Проверяем требуется ли сортировка в обратном порядке
	if *isReverseSorting {
		reverse(readResult)
	}

	//Проверяем требуется ли удаление повторяющихся элементов
	if *isRemoveRepeating {
		readResult = removeRepeats(readResult)
	}

	writeFileBytes(fileName, readResult)
}
