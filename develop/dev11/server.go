package main

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

//Run запускает сервер на указанном порту
//Принимает: строку с портом и http.Handler
//Возвращет ошибку от ListenAndServe
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, //1 MB
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

//Shutdown gracefully shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
