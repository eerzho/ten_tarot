package v1

import (
	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
	"gopkg.in/telebot.v3"
)

func NewHandler(l logger.Logger, bot *telebot.Bot, tgUserService *service.TGUser, tgMessageService *service.TGMessage) {
	// middleware
	mv := newMiddleware(l)

	// handler
	newCommand(l, bot, tgUserService)
	newMessage(l, mv, bot, tgMessageService, tgUserService)
}
