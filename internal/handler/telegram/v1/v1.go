package v1

import (
	"gopkg.in/telebot.v3"
)

func NewHandler(
	bot *telebot.Bot,
	tgUserService tgUserService,
	tgMessageService tgMessageService,
	tgButtonService tgButtonService,
	tgInvoiceService tgInvoiceService,
) {
	mv := newMiddleware(tgMessageService, tgButtonService, tgUserService)

	bot.Use(mv.setRIDAndLogDuration)

	newCommand(bot, tgUserService)
	newMessage(mv, bot, tgMessageService, tgUserService)
	newButton(bot, tgButtonService, tgInvoiceService)
	newPayment(bot, tgInvoiceService, tgUserService)
}
