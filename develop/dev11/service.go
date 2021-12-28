package main

import "time"

//Service тип, реализующий прослойку между handler и storage
type Service struct {
	storage *EventStore
}

//NewService конструктор, возвращающий ссылку на Serivce
func NewService(storage *EventStore) *Service {
	return &Service{storage}
}

//SaveEvent добавлет Event с переданными данными в хранилище
func (s *Service) SaveEvent(text string, date time.Time) error {
	return s.storage.Save(text, date)
}

//ChangeEvent изменяет Event из хранилища, заполняя новыми переданными данными
func (s *Service) ChangeEvent(id int, text string, date time.Time) (bool, error) {
	return s.storage.Change(id, text, date)
}

//GetTodays выдает слайс Event, относящихся к сегодняшнему дню
func (s *Service) GetTodays() ([]Event, error) {
	return s.storage.GetTodays()
}

//GetThisWeeks выдает слайс Event, относящихся к текущей неделе
func (s *Service) GetThisWeeks() ([]Event, error) {
	return s.storage.GetThisWeeks()
}

//GetThisMonths выдает слайс Event, относящихся к текущему месяцу
func (s *Service) GetThisMonths() ([]Event, error) {
	return s.storage.GetThisMonths()
}

//DeleteEvent удаляет Event из хранилища
func (s *Service) DeleteEvent(id int) (bool, error) {
	return s.storage.Delete(id)
}
