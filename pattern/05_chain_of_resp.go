package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного паттерна на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern

	Паттерн позволяет реализовать передачу запросов по последовательности
	обработчиков, где каждый обработчик решает может ли он обработать запрос
	и нужно ли передавать его дальше по цепи.
	Набор обработчиков может задаваться динамически, звено запроса всегда
	можно добавить, заменить или удалить.
	Цепь обработки запросов при этом может состоять лишь
	из одного обработчика и запрос может быть не обработан совсем.


	Плюсы:
	- Разнесение клиента и обработчиков, уменьшение их зависимости
	- Реализация принципа единственной ответственности

	Минусы:
	- Создание дополнительных объектов, усложнение кода
	- Запрос может быть не обработан

	Примером реализации паттерна может служить последовательность
	движения больного по больнице: он попадает в приемное отделение,
	затем к врачу, затем к кассиру на оплату приема.

*/

//Описание одного отделения больницы - один обработчик запроса
//Требует метод выполнения операции с пациентом - Execute
//и метод установки следующего обработчика запроса - SetNext
type Department interface {
	Execute(*Patient)
	SetNext(Department)
}

//Реализация отделения - приемное отделение, содержит поле -
//следующее отделение (следующий обработчик)
type Reseption struct {
	next Department
}

func (r *Reseption) Execute(p *Patient) {
	//Если пациент зарегистрирован просто вызовем следующий обработчик
	if p.RegistrationDone {
		fmt.Println("Patient already registered")
		r.next.Execute(p)
		return
	}

	//Если не зарегистрирован, то зарегистрируем и вызовем следующий обработчик
	fmt.Println("Registering patient")
	p.RegistrationDone = true

	r.next.Execute(p)
}

//Устанавливаем следующий обработчик запроса (следующее отделение)
func (r *Reseption) SetNext(d Department) {
	r.next = d
}

type Doctor struct {
	next Department
}

func (d *Doctor) Execute(p *Patient) {
	if p.DoctorDone {
		fmt.Println("Patient already visited Doctor")
		d.next.Execute(p)
		return
	}

	fmt.Println("Patient is visiting Doctor")
	p.DoctorDone = true

	d.next.Execute(p)
}

func (d *Doctor) SetNext(dep Department) {
	d.next = dep
}

type Cashier struct {
	next Department
}

func (c *Cashier) SetNext(d Department) {
	c.next = d
}

//Касса - итоговый объект обработки в цепочке,
//поэтому следующий объект - не вызываем
func (c *Cashier) Execute(p *Patient) {
	if p.CashierDone {
		fmt.Println("Patient already paid")
		return
	}
	fmt.Println("Patient is paying")
}

//Объект пациент с полями - статусами выполнения этапов
type Patient struct {
	RegistrationDone bool
	DoctorDone       bool
	CashierDone      bool
}

func main() {
	p := &Patient{} //создаем дефолтного пациента

	//создаем дефолтные объекты этапов обработки (этапов лечения пациента)
	r := &Reseption{}
	d := &Doctor{}
	c := &Cashier{}

	//Определяем последовательность этапов обработки
	r.SetNext(d) //посещение доктора - после регистратуры
	d.SetNext(c) //на кассу - после доктора

	r.Execute(p) //пациент приходит в приемное отделение
}
