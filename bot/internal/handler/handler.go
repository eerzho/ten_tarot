package handler

import (
	"bot/config"
	"bot/internal/repo/mongo_repo"
	"bot/internal/srv"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"
)

func SetUp(bot *telebot.Bot, cfg *config.Config, mng *mongo.Database, lg *slog.Logger) {
	// repo
	userRepo := mongo_repo.NewUser(lg, mng)
	invoiceRepo := mongo_repo.NewInvoice(lg, mng)
	messageRepo := mongo_repo.NewMessage(lg, mng)
	supportRequestRepo := mongo_repo.NewSupportRequest(lg, mng)

	// srv
	deckSrv := srv.NewDeck(lg)
	userSrv := srv.NewUser(lg, userRepo)
	tgCommandSrv := srv.NewTGCommand(lg)
	tgKeyboardSrv := srv.NewTGKeyboard(lg)
	invoiceSrv := srv.NewInvoice(lg, invoiceRepo, userSrv)
	tgInvoiceSrv := srv.NewTGInvoice(lg, invoiceSrv, userSrv)
	tarotSrv := srv.NewTarot(lg, cfg.GPT.Model, cfg.GPT.Token, cfg.GPT.Prompt)
	// tarotSrv := srv.NewTarotMock(lg)
	supportRequestSrv := srv.NewSupportRequest(lg, supportRequestRepo, userSrv)
	messageSrv := srv.NewMessage(lg, messageRepo, deckSrv, tarotSrv)

	mdw := newMiddleware(lg, userSrv, messageSrv, tgKeyboardSrv, 3)
	bot.Use(mdw.setRIDAndLogDuration, mdw.setUserAndContext)

	// handler
	newPayment(bot, lg, userSrv, invoiceSrv)
	newCommand(bot, lg, userSrv, tgCommandSrv)
	newButton(bot, lg, tgInvoiceSrv, tgKeyboardSrv)
	newMessage(bot, lg, mdw, userSrv, messageSrv, tgInvoiceSrv, supportRequestSrv)
}
