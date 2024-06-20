package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/config"
	"github.com/eerzho/ten_tarot/internal/app/http"
	"github.com/eerzho/ten_tarot/internal/repo/mongo_repo"
	"github.com/eerzho/ten_tarot/internal/service"
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

	l := logger.New(cfg.Level)

	// repo
	tgUserRepo := mongo_repo.NewTGUser(mg)
	tgMessageRepo := mongo_repo.NewTGMessage(mg)

	// service
	tgUserService := service.NewTGUser(l, tgUserRepo)
	cardService := service.NewCard()
	tarotService := service.NewTarot(l, cfg.GPT.Token, cfg.GPT.Prompt)
	tgMessageService := service.NewTGMessage(l, tgMessageRepo, tgUserService, cardService, tarotService)

	// handler
	httpServer := http.New(l, cfg, tgUserService, tgMessageService)

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
