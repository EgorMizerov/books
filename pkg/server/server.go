package server

import (
	"context"
	"net/http"
	"time"
)

//go:generate mockery --name=HttpServer
type HttpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Server struct {
	httpServer HttpServer
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       10 * time.Second,
		},
	}
}

func (s *Server) Listen() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
