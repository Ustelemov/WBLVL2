package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
	Объяснить применимость паттерна, его плюсы и минусы,
	а также реальные примеры использования данного паттерна на практике.
	https://en.wikipedia.org/wiki/State_pattern

	Паттерн позволяет менять объектам свое поведение, исходя из состояния
	в котором они находятся. Возможные состояния объекта выделяются
	в отдельные объекты-состояния и базовому объекту потребуется ссылаться
	только на текущее состояния и иметь возможность его изменять на другое.
	Так можно удобно добавлять новые состояния и изменять старые не
	внося изменения в код базового объекта. Большой условный оператор
	заменяется разными объектами-состояниями и ссылкой на них в базовом объекте.

	Плюс:
	- Уход от больших громоздких условных операторов
	- Код, относящийся к определенному состоянию находится в одном месте

	Минусы:
	- Может неоправдано усложнять код из-за введения дополнительных объектов,
	если состояний небольшое количество
	- Потребуется решать проблему доступа к
	 инкапсулированным данным контекста (базового класса)

	Для примера реализации рассмотрим умную колонку, которая
	может находится в трех состояниях:
	Проигрывается музыка
	Проигрывается подкаст
	Ничего не проигрывается - тишина

	И имеет три действия:
	Начать проигрывать музыку - можно корректно только из состоянии тишины
	Начать проигрывать подкаст - можно корректно только из состояния тишины
	И остановить вещание - сделать тишину: можно из музыки или подкаста
*/

//Определим базовый интерфейс состояния, требующий три метода
//соотвтетсвующих трем возможным действиям
type State interface {
	MusicOn()
	PodcastOn()
	StopBroadcasting()
}

//Реализуем умною колонку, хранящую своё имя и текущее состояние
type SmartSpeaker struct {
	name  string
	state State
}

//Реализуем метод установки текущего состояния
func (sp *SmartSpeaker) SetState(st State) {
	sp.state = st
}

//Реализуем методы действия для колонки, вызывая соответствующие
//методы у состояния
func (sp *SmartSpeaker) MusicOn() {
	sp.state.MusicOn()
}

func (sp *SmartSpeaker) PodcastOn() {
	sp.state.PodcastOn()
}

func (sp *SmartSpeaker) StopBroadcasting() {
	sp.state.StopBroadcasting()
}

//Реализуем конкретные состояния
//Если состояние нельзя выполнить, то выдаем соответствующий варнинг
type MusicState struct {
	smartSpeaker *SmartSpeaker
}

func (m *MusicState) MusicOn() {
	fmt.Println("Warning: Music is already playing")
}

func (m *MusicState) PodcastOn() {
	fmt.Println("Warning: Music is already playing")
}

func (m *MusicState) StopBroadcasting() {
	fmt.Println("OK: Music stopped")
	//Для примера в данной реализации будет передаваться новое
	//созданное состояние (объекты состояний можно было бы хранить
	// и в колонке, не создавая каждый раз новый)
	m.smartSpeaker.SetState(&StopState{
		smartSpeaker: m.smartSpeaker,
	})
}

type StopState struct {
	smartSpeaker *SmartSpeaker
}

func (m *StopState) MusicOn() {
	fmt.Println("OK: Music started playing")
	m.smartSpeaker.SetState(&MusicState{
		smartSpeaker: m.smartSpeaker,
	})
}

func (m *StopState) PodcastOn() {
	fmt.Println("OK: Podcast started playing")
	m.smartSpeaker.SetState(&PodcastState{
		smartSpeaker: m.smartSpeaker,
	})
}

func (m *StopState) StopBroadcasting() {
	fmt.Println("Warning: Stop is already set")
}

type PodcastState struct {
	smartSpeaker *SmartSpeaker
}

func (p *PodcastState) MusicOn() {
	fmt.Println("Warning: Podcast is already playing")
}

func (p *PodcastState) PodcastOn() {
	fmt.Println("Warning: Podcast is already playing")
}

func (p *PodcastState) StopBroadcasting() {
	fmt.Println("OK: Podcast stopped")
	p.smartSpeaker.SetState(&StopState{
		smartSpeaker: p.smartSpeaker,
	})
}

func main() {
	s := SmartSpeaker{
		name:  "Alice",
		state: nil,
	}
	s.SetState(&StopState{
		smartSpeaker: &s,
	})

	s.StopBroadcasting()
	s.MusicOn()
	s.MusicOn()
	s.PodcastOn()
	s.StopBroadcasting()
	s.StopBroadcasting()
	s.PodcastOn()
	s.StopBroadcasting()

}
