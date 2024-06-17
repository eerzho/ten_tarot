package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eerzho/event_manager/config"
	"github.com/eerzho/event_manager/internal/app/telegram"
	"github.com/eerzho/event_manager/internal/repo/mongo_repo"
	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/event_manager/pkg/mongo"
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

	l := logger.New(cfg.Level)

	// repo
	tgUserRepo := mongo_repo.NewTGUser(mg)
	tgMessageRepo := mongo_repo.NewTGMessage(mg)

	// service
	tgUserService := service.NewTGUser(l, tgUserRepo)
	eventService := service.NewEvent(l, cfg.GPT.Token, cfg.GPT.Prompt)
	googleCalendarService := service.NewGoogleCalendar(cfg.Google.CalendarURL)
	appleCalendar := service.NewAppleCalendar()
	tgMessageService := service.NewTGMessage(l, tgMessageRepo, tgUserService, eventService, googleCalendarService, appleCalendar)

	// handler
	telegramBot, err := telegram.New(l, cfg, tgUserService, tgMessageService)

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
