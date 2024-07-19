package v1

import (
	"context"
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/handler/http/v1/response"
	"github.com/eerzho/ten_tarot/internal/model"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	tgMessage struct {
		l                logger.Logger
		tgMessageService tgMessageService
	}

	tgMessageService interface {
		List(ctx context.Context, chatID string, page, count int) ([]model.TGMessage, int, error)
	}
)

func newTGMessage(l logger.Logger, router *gin.RouterGroup, tgMessageService tgMessageService) *tgMessage {
	t := tgMessage{
		l:                l,
		tgMessageService: tgMessageService,
	}

	router.GET("/tg-messages", t.list)

	return &t
}

// list -.
// @Summary Show messages
// @Description Show all messages list
// @Tags tg-messages
// @Accept json
// @Produce json
// @Param chat_id query string false "ChatID"
// @Param page query int false "Page"
// @Param count query int false "Count"
// @Success 200 {object} response.pagination{data=[]model.TGMessage}
// @Router /tg-messages [get]
func (t *tgMessage) list(ctx *gin.Context) {
	const op = "./internal/handler/http/v1/tg_message::list"

	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}

	count, _ := strconv.Atoi(ctx.Query("count"))
	if count == 0 {
		count = 10
	}

	messages, total, err := t.tgMessageService.List(ctx, ctx.Query("chat_id"), page, count)
	if err != nil {
		t.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
		response.Fail(ctx, err)
		return
	}

	response.Pagination(ctx, page, count, total, messages)
}
