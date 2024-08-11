package app

import (
	"bot/config"
	"bot/internal/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"
	"log"
	"log/slog"
)

type App struct {
	url string
	bot *telebot.Bot
}

func New(cfg *config.Config, mng *mongo.Database, lg *slog.Logger) (*App, error) {
	settings := telebot.Settings{
		Token: cfg.Telegram.Token,
		Poller: &telebot.Webhook{
			Listen: ":" + cfg.Telegram.Port,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: cfg.Telegram.Domain,
			},
		},
		OnError: nil,
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	handler.SetUp(bot, cfg, mng, lg)

	return &App{
		bot: bot,
		url: cfg.Telegram.Domain,
	}, nil
}

func (a *App) Run() {
	log.Printf("app: bot listening at %s", a.url)
	a.bot.Start()
}

func (a *App) Shutdown() {
	log.Print("app: bot shutting down")

	err := a.bot.RemoveWebhook()
	if err != nil {
		log.Printf("app: %v", err)
	}
}
