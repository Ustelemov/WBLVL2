package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

//initRouts инициализирует end-point'ы методами обработчиками
//Возвращает http.Handler
func (h *Handler) initRouts() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", h.createEvent)
	mux.HandleFunc("/update_event", h.updateEvent)
	mux.HandleFunc("/delete_event", h.deleteEvent)
	mux.HandleFunc("/events_for_day", h.eventsToday)
	mux.HandleFunc("/events_for_week", h.eventsThisWeek)
	mux.HandleFunc("/events_for_month", h.eventsThisMonth)

	handler := Logging(mux)

	return handler
}

//writeJsonMessage записывает в http.ResponseWriter сообщение
//в JSON формате с соответствующим хедером
//Принимает: статус код результата, строку сообщения и флаг ошибки
func writeJsonMessage(w http.ResponseWriter, status int, message string, isErr bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if isErr {
		data := struct {
			Message string `json:"Error"`
		}{Message: message}
		json.NewEncoder(w).Encode(data)
		return
	}

	data := struct {
		Message string `json:"Result"`
	}{Message: message}
	json.NewEncoder(w).Encode(data)
}

//writeJsonEvents записывает в http.ResponseWriter слайс Event
//в JSON формате с соответствующим хедером.
//Принимает: статус код результата, слайс Event.
func writeJsonEvents(w http.ResponseWriter, status int, events []Event) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	data := struct {
		Events []Event `json:"Result"`
	}{Events: events}
	json.NewEncoder(w).Encode(data)
}

//Logging middlware для логирования
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		log.Printf("%s %s", req.Method, req.RequestURI)
	})
}

//createEvent обработчик для POST /create_event
func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method POST at /create_event, got %v", r.Method), true)
		return
	}

	err := r.ParseForm()
	if err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Parse params error", true)
		return
	}

	text := r.Form.Get("text")
	dateStr := r.Form.Get("date")

	date, err := time.Parse("02-01-2006", dateStr)

	if dateStr == "" || err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Bad date", true)
		return
	}

	err = h.service.SaveEvent(text, date)

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, "Can't save event", true)
		return
	}

	writeJsonMessage(w, http.StatusOK, "Event saved", false)
}

//updateEvent обработчик для POST /update_event
func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method POST at /update_event, got %v", r.Method), true)
		return
	}

	err := r.ParseForm()
	if err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Parse params error", true)
		return
	}

	idStr := r.Form.Get("id")
	text := r.Form.Get("text")
	dateStr := r.Form.Get("date")

	id, err := strconv.Atoi(idStr)
	if id <= 0 || err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Bad id", true)
		return
	}

	date, err := time.Parse("02-01-2006", dateStr)

	//Если дата не была передана
	if dateStr == "" || err != nil {
		date = time.Time{}
	}

	isExists, err := h.service.ChangeEvent(id, text, date)

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, fmt.Sprintf("Can't change event: %s", err), true)
		return
	}

	if !isExists {
		writeJsonMessage(w, http.StatusBadRequest, fmt.Sprintf("Event with id %d doesn't exists", id), true)
		return
	}

	writeJsonMessage(w, http.StatusOK, "Event updated", false)
}

//deleteEvent обработчик для POST /delete_event
func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method POST at /delete_event, got %v", r.Method), true)
		return
	}

	err := r.ParseForm()
	if err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Parse params error", true)
		return
	}

	idStr := r.Form.Get("id")
	id, err := strconv.Atoi(idStr)
	if id <= 0 || err != nil {
		writeJsonMessage(w, http.StatusBadRequest, "Bad id", true)
		return
	}

	isExists, err := h.service.DeleteEvent(id)

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, fmt.Sprintf("Can't delete event: %s", err), true)
		return
	}

	if !isExists {
		writeJsonMessage(w, http.StatusBadRequest, fmt.Sprintf("Event with id %d doesn't exists", id), true)
		return

	}

	writeJsonMessage(w, http.StatusOK, "Event deleted", false)
}

//eventsToday обработчик для GET /events_for_day
func (h *Handler) eventsToday(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method Get at /events_for_day, got %v", r.Method), true)
		return
	}

	result, err := h.service.GetTodays()

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, fmt.Sprintf("Can't get events: %s", err), true)
		return
	}

	writeJsonEvents(w, http.StatusOK, result)
}

//eventsThisMonth обработчик для GET /events_for_month
func (h *Handler) eventsThisMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method Get at /events_for_month, got %v", r.Method), true)
		return
	}

	result, err := h.service.GetThisMonths()

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, fmt.Sprintf("Can't get events: %s", err), true)
		return
	}

	writeJsonEvents(w, http.StatusOK, result)
}

//eventsThisWeek обработчик для GET /events_for_week
func (h *Handler) eventsThisWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonMessage(w, http.StatusMethodNotAllowed, fmt.Sprintf("expect method Get at /events_for_week, got %v", r.Method), true)
		return
	}

	result, err := h.service.GetThisWeeks()

	if err != nil {
		writeJsonMessage(w, http.StatusServiceUnavailable, fmt.Sprintf("Can't get events: %s", err), true)
		return
	}

	writeJsonEvents(w, http.StatusOK, result)
}
