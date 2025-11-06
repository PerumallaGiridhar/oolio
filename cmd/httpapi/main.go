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

	index := index.NewPebbleIndex(cfg.PromoFiles)
	defer index.Close()

	validation.HTTPRequestValidatorInit(index)
	server := CreateServer(cfg.Server, routes.NewRouter())

	log.Printf("🚀 starting server on %s", cfg.Server.Addr)
	go server.Start()
	<-ctx.Done()
	log.Println("🛑 shutdown signal received")
	server.GraceFullShutdown(context.Background())
	log.Println("✅ graceful shutdown complete")
}
