package http

import (
	"autoscaler/internal/usecase"
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	handleRequestUseCase *usecase.HandleRequestUseCase
	server               *http.Server
}

func NewServer(handleRequestUseCase *usecase.HandleRequestUseCase) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
	}

	s := &Server{
		handleRequestUseCase: handleRequestUseCase,
		server:               server,
	}

	mux.HandleFunc("/", s.handleRequestUseCase.HandleRequest)
	return s
}

func (s *Server) Start(addr string) error {
	s.server.Addr = addr
	log.Printf("Starting server at %s\n", addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctxWithTimeout)
}
