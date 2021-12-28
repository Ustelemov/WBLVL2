package main

import (
	"fmt"
	"sync"
	"time"
)

//JSONTime дополнительный тип на основе time.Time
//для переопределения форматирования при маршалинге в JSON
type JSONTime time.Time

//MarshalJSON переопределяет маршалинг для типа JSONTime
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("02-01-2006"))
	return []byte(stamp), nil
}

//Event структура события
type Event struct {
	ID   int      `json:"ID"`
	Date JSONTime `json:"Date"`
	Text string   `json:"Text"`
}

//EventStore хранилище событий на основе map[int]Event
type EventStore struct {
	m      map[int]Event
	mutex  sync.RWMutex
	nextID int
}

//NewEventStore конструктор для EventStore.
//Возвращает: ссылку на созданный EventStore.
func NewEventStore() *EventStore {
	return &EventStore{m: make(map[int]Event, 0), nextID: 1}
}

//Save сохраняет новый Event в хранилище, присваивая ему порядковый id.
//Принимает строку текста события и время события.
//Возвращает ошибку сохранения.
//Конкурентно безопасный метод.
func (store *EventStore) Save(text string, date time.Time) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	event := Event{ID: store.nextID, Text: text, Date: JSONTime(date)}

	store.m[store.nextID] = event
	store.nextID = store.nextID + 1

	return nil
}

//Load получает Event из хранилища по id
//Принимает id искомого события.
//Возвращает объект события и флаг наличия\отсутствия элемента в хранилище.
//Конкурентно безопасный метод.
func (store *EventStore) Load(id int) (Event, bool) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	event, ok := store.m[id]
	return event, ok
}

//Change изменяет объект, находящийся в хранилище.
//Принимает все параметры cобытия: id, текст и время.
//Возвращает флаг наличия события ошибку изменения.
//Конкурентно безопасный метод.
func (store *EventStore) Change(id int, text string, date time.Time) (bool, error) {

	var el Event
	var ok bool
	if el, ok = store.Load(id); !ok {
		return false, nil
	}

	if date.IsZero() {
		date = time.Time(el.Date)
	}

	if text == "" {
		text = el.Text
	}

	event := Event{ID: id, Text: text, Date: JSONTime(date)}
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.m[id] = event

	return true, nil
}

//GetTodays возвращает все события, запланированные на сегодня.
//Возвращает: слайс событий и ошибку получения.
func (store *EventStore) GetTodays() ([]Event, error) {
	nowDay := time.Now()

	todayStart := time.Date(nowDay.Year(), nowDay.Month(), nowDay.Day(), 0, 0, 0, 0, nowDay.Location())
	todayEnd := todayStart.AddDate(0, 0, 1).Add(time.Nanosecond * -1)

	return store.getBetween(todayStart, todayEnd)
}

//GetThisWeeks возвращает все события, запланированные на текущую неделю.
//Возвращает: слайс событий и ошибку получения.
func (store *EventStore) GetThisWeeks() ([]Event, error) {
	nowDay := time.Now()
	firstWeekDay := nowDay

	for firstWeekDay.Weekday() != time.Monday {
		firstWeekDay = firstWeekDay.AddDate(0, 0, -1)
	}
	firstWeekDay = time.Date(firstWeekDay.Year(), firstWeekDay.Month(),
		firstWeekDay.Day(), 0, 0, 0, 0, firstWeekDay.Location())
	lastWeekDay := firstWeekDay.AddDate(0, 0, 7).Add(time.Nanosecond * -1)

	return store.getBetween(firstWeekDay, lastWeekDay)
}

//GetThisMonths возвращает все события, запланированные на текущий месяц.
//Возвращает: слайс событий и ошибку получения.
func (store *EventStore) GetThisMonths() ([]Event, error) {
	nowDay := time.Now()

	monthStart := time.Date(nowDay.Year(), nowDay.Month(), 1, 0, 0, 0, 0, nowDay.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	return store.getBetween(monthStart, monthEnd)
}

//getBetween возвращает все события, запланированные в промежутке между временем cтарта и окончания
//Принимает: время старта и время окончания.
//Возвращает: слайс событий и ошибку получения.
//Конкурентно безопасный метод.
func (store *EventStore) getBetween(start time.Time, end time.Time) ([]Event, error) {
	result := make([]Event, 0)
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	for _, v := range store.m {
		store.mutex.RUnlock()

		if inTimeSpan(start, end, time.Time(v.Date)) {
			result = append(result, v)
		}
		store.mutex.RLock()
	}
	return result, nil
}

//inTimeSpan проверяет нахождение времени в заданном промежутке.
//Принимает: время старта промежутка, время окончания промежутка и проверяемое время.
//Возвращает: булевский результат нахождения.
func inTimeSpan(start, end, check time.Time) bool {
	if start.After(end) {
		start, end = end, start
	}

	return !check.Before(start) && !check.After(end)
}

//Delete удаляет объект из хранилища по id.
//Принимет: id удаляемого объекта.
//Возвращает: булевский результат нахождения объекта и ошибку удаления.
func (store *EventStore) Delete(id int) (bool, error) {
	if _, ok := store.Load(id); !ok {
		return false, nil
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.m, id)

	return true, nil
}
