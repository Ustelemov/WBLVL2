package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного паттерна на практике.
	https://en.wikipedia.org/wiki/Facade_pattern

	Патттерн может использоваться, когда нам нужно работать с некоторым набором объектов сложноустроенной
	подсистемы (например, библиотеки), где каждый объект требуется инициализировать определенным образом,
	отслеживать зависимости, вызывать методы в определенном порядке и пр.
	Это усложняет работу и если клиенту необходима какая-то небольшая функциональность, контроль всего этого
	сильно усложнит работу для него.
	Поэтому мы можем скрыть сложные детали работы с подсистемой за фасадом, предоставив пользователю лишь
	простой интерфейс для выполнения желаемой функциональности.
	Так, мы не только предоставляем простой и удобный способ работы для клиента,
	но и минимизируем его зависимость от используемой подсистемы.

	Плюсами такого подхода является:
	- Упрощение работы клиента с подсистемой - меньше кода, меньше ошибок, быстрее разработка.
	- Уменьшении зависимости от подсистемы - проще внести изменения, проще тестировать.
	- Упрощение внешней документации - упрощение работы с подсистемой для клиента - проще клиентская документация

	Минусами подхода является:
	- Требуется дополнительная реализация необходимых интерфейсов - дополнительная разработка.
	- Нужно хорошо продумать реализуемый набор интерфейсов для клиента, чтобы вся функциональность, ему
	необходимая, была у него доступна (при доработках подсистемы нужно поддерживать и фасад).
	Если необходимой функциональности нет - клиенту придется реализовывать её самому
	(в обход паттерна), а нам затем реализовывать её в рамках паттерна.

	Примером использования паттерна может служить реализация компьютера.
	В данном случае - кнопка включения компьютера является фасадным интерфейсом для различных операций,
	происходимых для запуска компьютера, к примеру: запуск BIOS, запуск Операционной системы.
	Причем вызов методов запуска должен происходить именно в таком порядке. Мы скрываем детали за
	методом Start структуры ComputerFacade.

	Посути, создание некоторого метода-обертки, например: открытие соединения к БД, получение данных,
	закрытие соединения, агрегация данных, возвращение результата - также является реализацией такого паттерна.
	Мы предоставляем простой интерфейс, скрываем некоторую усложненную логику, требующую определенного порядка действий.
*/

//Объект BIOS
type BIOS struct{}

//Объект Операционная система
type OperationSystem struct{}

//Метод запуска (настройки) BIOS
func (bios *BIOS) Setup() {
	fmt.Println("Starting BIOS...")
}

//Метод запуска (настройки) Операционной системы
func (os *OperationSystem) Setup() {
	fmt.Println("Starting Operation system...")
}

//Объект Фасада, в данном случае - сам компьютер или точнее
//системный блок, на котором будет нажата кнопка запуска (метод Start)
//В нем находятся поля для обоих объектов, методы которых будем обворачивать
type ComputerFacade struct {
	bios *BIOS
	os   *OperationSystem
}

//Конструктор без параметров, просто создаем пустые объекты в полях
func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{os: &OperationSystem{}, bios: &BIOS{}}
}

//Основной метод Фасад, представляющий собой нажатие кнопки на системном
//блоке, запускаем BIOS, затем операционную систему - порядок важен!
func (cf *ComputerFacade) Start() {
	cf.bios.Setup()
	cf.os.Setup()
	fmt.Println("Computer is started")
}

func main() {
	cf := NewComputerFacade() //создаем объект Фасада
	cf.Start()                //нажимаем кнопку старт на системном блоке
}
