package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного паттерна на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern

	Паттерн используется, когда нужно добавить некоторую функциональность
	к объектам некоторой иерархии, не изменяя сами объекты.

	Плюсы:
	- Позволяет полагаться на стабильность объектов, так как не требует внесений
	изменений в них
	- Схожие операции для разных объектов добавляются и располагаются в одном классе -
	упрощается навигация, внесение изменений и тестирование

	Минусы:
	- Может нарушаться инкапсуляция (если посетителю потребуется доступ к приватным полям)

	Примером реализации паттерна может послужить добавление разных принтеров: в JSON и XML формате
	для фигур Квадрат, Треугольник и Круг, реализуемых интерфейс Фигура.
*/

//Базовый интерфейс Фигура - корень иерархии, требующий метод принятия визитора
type Shape interface {
	Accept(Visitor)
}

//Объект квадрат с полем - стороной
type Square struct {
	side int
}

//Реализация метода принятия визитора
func (s *Square) Accept(v Visitor) {
	v.VisitSquare(*s)
}

//Объект круг с полем - радиусом
type Circle struct {
	radius int
}

//Реализация метода принятия визитора
func (c *Circle) Accept(v Visitor) {
	v.VisitCircle(*c)
}

//Объект прямоугольник с полями - сторонами
type Rectangle struct {
	a, b int
}

//Реализация метода принятия визитора
func (r *Rectangle) Accept(v Visitor) {
	v.VisitRectangle(*r)
}

//Интерфейс визитор, требующий реализации методов посещения для каждого
//объекта из иерархии
type Visitor interface {
	VisitSquare(Square)
	VisitCircle(Circle)
	VisitRectangle(Rectangle)
}

//Реализация интерфейса визитора для печати фигур в XML формате
type XMLPrinter struct{}

func (x *XMLPrinter) VisitSquare(Square) {
	fmt.Printf("<xml>Square</xml>\n")
}

func (x *XMLPrinter) VisitCircle(Circle) {
	fmt.Printf("<xml>Circle</xml>\n")
}

func (x *XMLPrinter) VisitRectangle(Rectangle) {
	fmt.Printf("<xml>Rectangle</xml>\n")
}

//Реализация интерфейса визитора для печати фигур в JSON формате
type JSONPrinter struct{}

func (j *JSONPrinter) VisitSquare(Square) {
	fmt.Println(`{"Square"}`)
}

func (j *JSONPrinter) VisitCircle(Circle) {
	fmt.Println(`{"Circle"}`)
}

func (j *JSONPrinter) VisitRectangle(Rectangle) {
	fmt.Println(`{"Rectangle"}`)
}

func main() {
	//создаем объекты иерархии
	r := Rectangle{}
	c := Circle{}
	s := Square{}

	xml := XMLPrinter{} //создаем объект реализующий визитор
	//вызываем метод приема визитора для каждого объета иерархии
	r.Accept(&xml)
	c.Accept(&xml)
	s.Accept(&xml)

	json := JSONPrinter{} //создаем объект реализующий визитор
	//вызываем метод приема визитора для каждого объета иерархии
	r.Accept(&json)
	c.Accept(&json)
	s.Accept(&json)
}
