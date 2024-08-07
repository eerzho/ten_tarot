package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eerzho/ten_tarot/config"
	"github.com/eerzho/ten_tarot/internal/app/http"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/eerzho/ten_tarot/pkg/mongo"
)

func main() {
	const op = "./cmd/http_server::main"

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	mg, err := mongo.New(cfg.Mongo.URL, cfg.Mongo.DB)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
	defer mg.Close()

	logger.SetUpLogger(cfg.Log.Level, cfg.Log.Type)

	httpServer := http.New(cfg, mg)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		httpServer.Run()
	}()

	log.Printf("%s: http server started", op)
	<-stopChan

	log.Printf("%s: shutting down", op)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)

	log.Printf("%s: http server stopped", op)
}
