package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

/*
Необходимо реализовать свою собственную UNIX-шелл-утилиту с
поддержкой ряда простейших команд:
- cd <args> - смена директории (в качестве аргумента могут
быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте
аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в
формате *такой-то формат*
Так же требуется поддерживать функционал fork/exec-команд
Дополнительно необходимо поддерживать конвейер на пайпах
(linux pipes, пример cmd1 | cmd2 | .... | cmdN).
*Шелл — это обычная консольная программа, которая будучи
запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись
ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный
сеанс поддерживается до тех пор, пока не будет введена
команда выхода (например \quit).*/

//proc структура процесса, используемая для сортировки по pid
//name строка, имя процесса
//pid целое неотрицательное число (id процесса)
type proc struct {
	name string
	pid  int
}

//cmd базовый интерфейс команды, требующий метод выполнения exec.
//Exec принимает массив строк аргументов команды, буфер байт для ввода данных
//от предыдущей команды пайплайна, флаг пайплайна (в команде определяется
//использование\неиспользование в пайплайне).
//Возвращает буфер байт с результатом выполнения команды и ошибку выполнения.
type cmd interface {
	exec([]string, *bytes.Buffer, bool) (*bytes.Buffer, error)
}

//структуры команд, реализующие cmd
type cdCMD struct{}
type echoCMD struct{}
type pwdCMD struct{}
type psCMD struct{}
type killCMD struct{}
type execCMD struct{}
type forkCMD struct{}
type exitCMD struct{}

//команда CD
func (cmd *cdCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	//cd в пайплайне не выполняет действия
	if !chain {
		//требуется путь изменения директории
		if len(args) > 0 {
			return nil, os.Chdir(args[0])
		}
	}
	return nil, nil
}

//команда Echo
func (cmd *echoCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var err error

	//c аргументом - выводит аргумент, без аргумента - пустую строку
	if len(args) > 0 {
		_, err = buf.WriteString(args[0] + "\n")
	} else {
		_, err = buf.WriteString("\n")
	}

	if err != nil {
		return nil, err
	}

	return &buf, nil
}

//команда Exec
func (cmd *execCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	c := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer

	//если от предыдущей команды есть данные для ввода в текущую команду
	//(реализация linux pipe)
	if in != nil {
		c.Stdin = bytes.NewReader((*in).Bytes())
	}

	//вывод команды в буфер байт
	c.Stdout = &out
	c.Stderr = os.Stderr

	c.Run()

	os.Exit(0)

	return &out, nil
}

//команда Fork
func (cmd *forkCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	c := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer

	//если от предыдущей команды есть данные для ввода в текущую команду
	//(реализация linux pipe)
	if in != nil {
		c.Stdin = bytes.NewReader((*in).Bytes())
	}

	//вывод команды в буфер байт
	c.Stdout = &out
	c.Stderr = os.Stderr
	c.Run()

	return &out, nil
}

//команда Pwd
func (cmd *pwdCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	var out bytes.Buffer

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	out.WriteString(pwd + "\n")
	return &out, nil
}

//команда Ps
func (cmd *psCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	var out bytes.Buffer

	//получаем все процессы
	matches, err := filepath.Glob("/proc/*/exe")
	if err != nil {
		return nil, err
	}

	procs := make([]proc, 0)

	for _, file := range matches {
		target, _ := os.Readlink(file) //получаем ссылку с именем процесса
		if len(target) > 0 {
			name := filepath.Base(target) //получаем имя процесса из пути
			pid, err := strconv.Atoi(strings.Split(file, "/")[2])

			if err == nil {
				procs = append(procs, proc{name, pid})
			}

		}
	}

	//сортируем по pid
	sort.Slice(procs, func(i, j int) bool {
		return procs[i].pid < procs[j].pid
	})

	//шапка столбцов
	out.WriteString("PID Name")

	for _, v := range procs {
		out.WriteString(fmt.Sprintf("%d %s\n", v.pid, v.name))
	}
	return &out, nil
}

//команда kill
func (cmd *killCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	if len(args) > 0 {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}
		syscall.Kill(pid, syscall.SIGKILL)
	}
	return nil, nil
}

//команда exit
func (cmd *exitCMD) exec(args []string, in *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	if !chain { //в пайплайне не выполняет действий
		os.Exit(0)
	}
	return nil, nil
}

//execCommand выполняет команду, разбивая входную строку команды по пробелам
//и используя нулевой элемент разбиения как название команды, а остальное - как параметры.
//Принимает строку команды, буфера с результатом предудыщей команды,
//флаг цепочки команд.
//Возвращает буфер с результатом и ошибку выполнения.
func execCommand(line string, prev *bytes.Buffer, chain bool) (*bytes.Buffer, error) {
	splited := strings.Split(line, " ")
	args := make([]string, 0)

	//удалим лишние пустые строки
	for i := range splited {
		if splited[i] != "" {
			args = append(args, splited[i])
		}
	}

	var cmd cmd

	switch args[0] {
	case "cd":
		cmd = &cdCMD{}
	case "pwd":
		cmd = &pwdCMD{}
	case "echo":
		cmd = &echoCMD{}
	case "ps":
		cmd = &psCMD{}
	case "kill":
		cmd = &killCMD{}
	case "exec":
		cmd = &execCMD{}
	case "fork":
		cmd = &forkCMD{}
	case "exit":
		cmd = &exitCMD{}
	}

	if cmd != nil {
		return cmd.exec(args[1:], prev, chain)
	}

	return nil, nil
}

//execCommands выполняет команды, разбивая входную строку команд
//по символу пайплайна "|" и выполняя все команды слева на право.
//Принимает строку команд.
//Возвращает ошибку выполнения.
func execCommands(line string) error {
	commands := strings.Split(line, "|")
	chain := len(commands) > 1

	var prev *bytes.Buffer
	var err error

	for _, v := range commands {
		prev, err = execCommand(v, prev, chain)
		if err != nil {
			return err
		}
	}

	if prev == nil || len(prev.String()) == 0 {
		fmt.Fprint(os.Stdout, "")
	} else {
		fmt.Fprint(os.Stdout, prev.String())
	}

	return err
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	//бесконечно пока не будет введено "exit"
	for {
		pwd, err := os.Getwd() //получаем текущий путь для вывода

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s$ ", pwd)
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		//Удалим символ перевода строки
		commands := strings.TrimSuffix(input, "\n")

		err = execCommands(commands)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
