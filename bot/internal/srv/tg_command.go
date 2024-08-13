package srv

import (
	"bot/internal/def"
	"context"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type TGCommand struct {
	lg *slog.Logger
}

func NewTGCommand(
	lg *slog.Logger,
) *TGCommand {
	return &TGCommand{
		lg: lg,
	}
}

func (tc *TGCommand) GetCommands(ctx context.Context) []telebot.Command {
	const op = "srv.TGCommand.SetCommands"
	tc.lg.Debug(op)

	commands := []telebot.Command{
		{Text: def.TGDonateCommand, Description: "Пожертвовать в развитие проекта"},
		{Text: def.TGSupportCommand, Description: "Связаться с разработчиком"},
		{Text: def.TGCancelCommand, Description: "Отменить активную команду"},
	}

	return commands
}
