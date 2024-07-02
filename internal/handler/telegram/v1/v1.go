package v1

import (
	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/service"
	"gopkg.in/telebot.v3"
)

func NewHandler(l logger.Logger, bot *telebot.Bot, tgUserService *service.TGUser, tgMessageService *service.TGMessage) {
	// middleware
	mv := newMiddleware(l, tgMessageService)

	// handler
	newCommand(l, bot, tgUserService)
	newMessage(l, mv, bot, tgMessageService, tgUserService)
}
