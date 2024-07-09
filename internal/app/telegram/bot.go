package telegram

import (
	"fmt"
	"log"
	"strings"

	"github.com/eerzho/event_manager/pkg/crypter"
	"github.com/eerzho/event_manager/pkg/mongo"
	"github.com/eerzho/ten_tarot/config"
	v1 "github.com/eerzho/ten_tarot/internal/handler/telegram/v1"
	"github.com/eerzho/ten_tarot/internal/repo/mongo_repo"
	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	url string
	bot *telebot.Bot
}

func New(l logger.Logger, cfg *config.Config, mg *mongo.Mongo, c *crypter.Crypter) (*Bot, error) {
	url := fmt.Sprintf("%s/ten-tarot/wb", strings.Trim(cfg.Telegram.Domain, "/"))
	settings := telebot.Settings{
		Token: cfg.Telegram.Token,
		Poller: &telebot.Webhook{
			Listen: ":" + cfg.Telegram.Port,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: url,
			},
		},
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("./internal/app/bot::New: %w", err)
	}

	// repo
	tgUserRepo := mongo_repo.NewTGUser(mg)
	tgMessageRepo := mongo_repo.NewTGMessage(c, mg)

	// service
	tgUserService := service.NewTGUser(tgUserRepo)
	cardService := service.NewCard()
	// todo
	//tarotService := service.NewTarotMock()
	tarotService := service.NewTarot(cfg.Model, cfg.GPT.Token, cfg.GPT.Prompt)
	tgMessageService := service.NewTGMessage(tgMessageRepo, cardService, tarotService)

	v1.NewHandler(l, bot, tgUserService, tgMessageService)

	return &Bot{
		bot: bot,
		url: url,
	}, nil
}

func (t *Bot) Run() {
	const op = "./internal/app/telegram/bot::Run"

	log.Printf("%s: telegram bot listening at %s", op, t.url)
	t.bot.Start()
}

func (t *Bot) Shutdown() {
	const op = "./internal/app/telegram/bot::Shutdown"

	log.Printf("%s: telegram bot shutting down", op)
	err := t.bot.RemoveWebhook()
	if err != nil {
		log.Printf("%s: %v", op, err)
	}
}
