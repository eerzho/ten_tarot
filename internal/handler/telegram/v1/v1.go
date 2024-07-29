package v1

import (
	"gopkg.in/telebot.v3"
)

func NewHandler(
	bot *telebot.Bot,
	tgUserService tgUserService,
	tgMessageService tgMessageService,
	tgKeyboardService tgKeyboardService,
	tgInvoiceService tgInvoiceService,
) {
	mv := newMiddleware(tgMessageService, tgKeyboardService, tgUserService)

	bot.Use(mv.setRIDAndLogDuration, mv.setUserAndContext)

	newCommand(bot, tgUserService)
	newMessage(mv, bot, tgMessageService, tgUserService)
	newButton(bot, tgKeyboardService, tgInvoiceService)
	newPayment(bot, tgInvoiceService, tgUserService)
}
