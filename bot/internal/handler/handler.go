package handler

import (
	"bot/config"
	"bot/internal/repo/mongo_repo"
	"bot/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"
	"log/slog"
)

func SetUp(bot *telebot.Bot, cfg *config.Config, mng *mongo.Database, lg *slog.Logger) {
	// repo
	tgUserRepo := mongo_repo.NewTGUser(lg, mng)
	tgMessageRepo := mongo_repo.NewTGMessage(lg, mng)
	tgInvoiceRepo := mongo_repo.NewTGInvoice(lg, mng)
	supportRequestRepo := mongo_repo.NewSupportRequest(lg, mng)

	// service
	deckService := service.NewDeck(lg)
	tgKeyboardService := service.NewTGKeyboard(lg)
	tgUserService := service.NewTGUser(lg, tgUserRepo)
	tgInvoiceService := service.NewTGInvoice(lg, tgInvoiceRepo, tgUserService)
	// tarotService := service.NewTarot(lg, cfg.GPT.Model, cfg.GPT.Token, cfg.GPT.Prompt)
	tarotService := service.NewTarotMock(lg)
	tgMessageService := service.NewTGMessage(lg, tgMessageRepo, deckService, tarotService)
	supportRequestService := service.NewSupportRequest(lg, supportRequestRepo, tgUserService)

	mdw := newMiddleware(lg, tgMessageService, tgKeyboardService, tgUserService)

	bot.Use(mdw.setRIDAndLogDuration, mdw.setUser)

	newCommand(bot, lg, tgUserService)
	newPayment(bot, lg, tgInvoiceService, tgUserService)
	newButton(bot, lg, tgKeyboardService, tgInvoiceService)
	newMessage(bot, lg, mdw, tgMessageService, tgUserService, tgInvoiceService, supportRequestService)
}
