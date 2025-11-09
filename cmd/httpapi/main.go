package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/PerumallaGiridhar/oolio/internal/config"
	"github.com/PerumallaGiridhar/oolio/internal/index"
	"github.com/PerumallaGiridhar/oolio/internal/routes"
	"github.com/PerumallaGiridhar/oolio/internal/validation"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	log.Printf("Initializing pebble store")
	index, err := index.NewPebbleIndex(cfg.PromoFiles)
	if err != nil {
		log.Fatalf("initializing pebble index: %v", err)
	}
	defer index.Close()

	if err := validation.HTTPRequestValidatorInit(index); err != nil {
		log.Fatalf("initializing HTTP request validator: %v", err)
	}

	server := CreateServer(cfg.Server, routes.NewRouter())

	log.Printf("🚀 starting server on %s", cfg.Server.Addr)
	go server.Start()
	<-ctx.Done()
	log.Println("🛑 shutdown signal received")
	server.GraceFullShutdown(ctx)
	log.Println("✅ graceful shutdown complete")
}
