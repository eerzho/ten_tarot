package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eerzho/event_manager/pkg/crypter"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/config"
	"github.com/eerzho/ten_tarot/internal/app/telegram"
	"github.com/eerzho/ten_tarot/pkg/logger"
)

func main() {
	const op = "./cmd/telegram_bot::main"

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	mg, err := mongo.New(cfg.Mongo.URL, cfg.Mongo.DB)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}
	defer mg.Close()

	l := logger.New(cfg.Log.Level)
	c := crypter.New(cfg.Crypter.Key)

	telegramBot, err := telegram.New(l, cfg, mg, c)

	if err != nil {
		log.Fatalf("%s: %s", op, err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		telegramBot.Run()
	}()

	log.Printf("%s: telegram bot started", op)
	<-stopChan

	log.Printf("%s: shutting down", op)

	telegramBot.Shutdown()

	log.Printf("%s: telegram bot stopped", op)
}
