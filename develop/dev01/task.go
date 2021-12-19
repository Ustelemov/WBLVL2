package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*
=== Базовая задача ===
Создать программу печатающую точное время с использованием NTP библиотеки.
Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время/точное время с использованием этой библиотеки.
Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

//getSystemUnixTime возвращает строку, содержащую текущее системное время
//в time.UnixDate формате
func getSystemUnixTime() string {
	return time.Now().Local().Format(time.UnixDate)
}

//getNetUnixTime возвращает строку, содержащую текущее сетевое время
//в time.UnixDate формате и ошибку получения сетевого времени
//Параметр server - строка, задающая адрес сервера для получения
//сетевого времени
func getNetUnixTime(server string) (string, error) {
	ntpTime, err := ntp.Time(server)

	if err != nil {
		return "", err
	}

	ntpTimeFormatted := ntpTime.Format(time.UnixDate)

	return ntpTimeFormatted, nil

}

const server = "time.apple.com"

func main() {

	sysTime := getSystemUnixTime()

	fmt.Printf("System time (Unix Date): %v\n", sysTime)

	netTime, err := getNetUnixTime(server)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Network time (Unix Date): %v\n", netTime)
}
