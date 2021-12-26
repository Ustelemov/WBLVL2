package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

/*
=== Утилита wget ===
Реализовать утилиту wget с возможностью скачивать сайты целиком
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//Скачивание сайта всего сайт по ссылке. Сохранение в текущей дериктории,
//корневая папка - домен сайта
func main() {

	//Ввод URL-сайта
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("$wget ")
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSuffix(input, "\n")

	//Парсинг URL (проверка корректности)
	u, err := url.ParseRequestURI(input)
	if err != nil {
		log.Fatal(err)
	}

	hostname := u.Hostname()

	//создадим корневую папку, куда будем скачивать сайт
	os.Mkdir(hostname, os.ModePerm)

	//создадим мапу для проверки обрабатывали уже эту ссылку или нет
	seen := make(map[string]bool)

	//опишем все доступные url в регулярном выражении
	//может включать все поддомены и http c https
	reg, err := regexp.Compile("https?://([a-z0-9]+[.])*" + hostname)

	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector(
		colly.URLFilters(reg),
	)

	//берем каждый a тег, получаем ссылку из атрибута href
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		absLink := e.Request.AbsoluteURL(link)

		//если ещё не посещали ссылку - посетим
		if !seen[absLink] {
			c.Visit(absLink)
			seen[absLink] = true
		}
	})

	//по ответу - записываем в файл
	c.OnResponse(func(r *colly.Response) {

		p := r.Request.URL.Path  //путь после домена
		fullPath := hostname + p //полный путь

		//Если нет таких директорий - создадим
		if _, err := os.Stat(fullPath); err != nil {
			os.MkdirAll(fullPath, os.ModePerm)
		}

		//если конечная ссылка, значит сохраняем как index.html
		if path.Ext(p) == "" {
			//проверим есть ли в конце "/"
			if fullPath[len(fullPath)-1] != '/' {
				fullPath += "/"
			}
			fullPath += "index.html"
		}

		fmt.Printf("Saving: %s\n", p)
		r.Save(fullPath)
	})

	//Начнем с переданного сайта обход
	if err := c.Visit(u.String()); err != nil {
		log.Fatal(err)
	}
	c.Wait()
	fmt.Printf("Fully downloaded %s\n", u.String())
}
