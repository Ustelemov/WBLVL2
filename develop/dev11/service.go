package main

import "time"

type Service struct {
	storage *EventStore
}

func NewService(storage *EventStore) *Service {
	return &Service{storage}
}

func (s *Service) SaveEvent(text string, date time.Time) error {
	return s.storage.Save(text, date)
}

func (s *Service) ChangeEvent(id int, text string, date time.Time) (bool, error) {
	return s.storage.Change(id, text, date)
}

func (s *Service) GetTodays() ([]Event, error) {
	return s.storage.GetTodays()
}

func (s *Service) GetThisWeeks() ([]Event, error) {
	return s.storage.GetThisWeeks()
}

func (s *Service) GetThisMonths() ([]Event, error) {
	return s.storage.GetThisMonths()
}

func (s *Service) DeleteEvent(id int) (bool, error) {
	return s.storage.Delete(id)
}
