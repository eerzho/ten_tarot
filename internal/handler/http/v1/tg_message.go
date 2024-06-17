package v1

import (
	"fmt"
	"strconv"

	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/gin-gonic/gin"
)

type tgMessage struct {
	l                logger.Logger
	tgMessageService *service.TGMessage
}

func newTGMessage(l logger.Logger, router *gin.RouterGroup, tgMessageService *service.TGMessage) *tgMessage {
	t := tgMessage{
		l:                l,
		tgMessageService: tgMessageService,
	}

	router.GET("/tg-messages", t.all)

	return &t
}

// @Summary Show messages
// @Description Show all messages list
// @Tags tg-messages
// @Accept json
// @Produce json
// @Param chat_id query string false "ChatID"
// @Param page query int false "Page"
// @Param count query int false "Count"
// @Success 200 {object} successResponse{data=[]entity.TGMessage}
// @Failure 500 {object} errorResponse
// @Router /tg-messages [get]
func (t *tgMessage) all(ctx *gin.Context) {
	const op = "./internal/handler/http/v1/tg_message::all"

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		count = 0
	}

	messages, err := t.tgMessageService.All(ctx, ctx.Query("chatID"), page, count)
	if err != nil {
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		errorRsp(ctx, err)
		return
	}

	successRsp(ctx, messages)
	return
}
