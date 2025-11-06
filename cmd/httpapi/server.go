package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PerumallaGiridhar/oolio/internal/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	ChiRouter *chi.Mux
	server    *http.Server
	config    config.ServerConfig
}

func (s *Server) Start() {
	s.server = &http.Server{
		Addr:              s.config.Addr,
		Handler:           s.ChiRouter,
		ReadHeaderTimeout: time.Duration(s.config.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.config.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.config.IdleTimeout) * time.Second,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func (s *Server) GraceFullShutdown(ctx context.Context) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	log.Println("Shutting down server gracefully...")
	if err := s.server.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}

func CreateServer(config config.ServerConfig, router *chi.Mux) *Server {
	return &Server{
		ChiRouter: router,
		config:    config,
	}
}
