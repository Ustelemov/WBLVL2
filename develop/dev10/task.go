package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

/*
Реализовать простейший telnet-клиент.
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123
Требования:
1. Программа должна подключаться к указанному хосту (ip или
доменное имя + порт) по протоколу TCP. После подключения
STDIN программы должен записываться в сокет, а данные
полученные и сокета должны выводиться в STDOUT
2. Опционально в программу можно передать таймаут на
подключение к серверу (через аргумент --timeout, по
умолчанию 10s)
3. При нажатии Ctrl+D программа должна закрывать сокет и
завершаться. Если сокет закрывается со стороны сервера,
программа должна также завершаться. При подключении к
несуществующему сервер, программа должна завершаться
через timeout
*/

func main() {
	timeoutStr := flag.String("timeout", "10s", "Timeout for connection")
	flag.Parse()

	timeout, err := time.ParseDuration(*timeoutStr)
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() != 2 {
		log.Fatal(fmt.Errorf("bad count of args: %d", flag.NArg()))
	}

	args := flag.Args()

	ip := args[0]
	port := args[1]
	addr := ip + ":" + port

	fmt.Println("Connecting...")
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")

	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Command to send: ")
		scanned := scanner.Scan()

		//Если EOF = CTRL+D
		if !scanned && scanner.Err() == nil {
			break
		}

		text := scanner.Text()

		//Пишем в сокет
		fmt.Fprintf(conn, "%s\n", text)

		//Читаем в буфер (при необходимости расширим \ заменим на несколько с добавлением)
		b := make([]byte, 512)
		n, err := conn.Read(b)

		//Если соединение закрыто
		if n == 0 || err == io.EOF {
			fmt.Println("Connection closed")
			break
		}

		//Если другая ошибка
		if err != nil {
			fmt.Printf("Error while read: %s\n", err)
			break
		}

		fmt.Fprintln(os.Stdout, string(b))
	}
	fmt.Println() //отступ в конце для красоты
}
