package app

import (
	"bot/config"
	"bot/internal/handler"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"
)

type App struct {
	bot *telebot.Bot
}

func New(cfg *config.Config, mng *mongo.Database, lg *slog.Logger) (*App, error) {
	settings := telebot.Settings{
		Token: cfg.Telegram.Token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
		OnError: nil,
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	handler.SetUp(bot, cfg, mng, lg)

	return &App{bot: bot}, nil
}

func (a *App) Run() {
	log.Println("app: bot started with long polling")
	if err := a.bot.RemoveWebhook(); err != nil {
		log.Fatalf("app: %v", err)
	}
	a.bot.Start()
}

func (a *App) Shutdown() {
	log.Println("app: bot shutting down")
}
