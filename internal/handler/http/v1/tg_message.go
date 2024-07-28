package v1

import (
	"strconv"

	"github.com/eerzho/ten_tarot/internal/handler/http/v1/response"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	tgMessage struct {
		tgMessageService tgMessageService
	}
)

func newTGMessage(router *gin.RouterGroup, tgMessageService tgMessageService) *tgMessage {
	t := tgMessage{
		tgMessageService: tgMessageService,
	}

	router.GET("/tg-messages", t.getList)

	return &t
}

// @Summary Show messages
// @Description Show messages list
// @Tags tg-messages
// @Accept json
// @Produce json
// @Param chat_id query string false "ChatID"
// @Param page query int false "Page"
// @Param count query int false "Count"
// @Success 200 {object} response.pagination{data=[]model.TGMessage}
// @Router /tg-messages [get]
func (t *tgMessage) getList(ctx *gin.Context) {
	const op = "handler.http.v1.tgMessage.getList"

	// todo перемести эту логику в service
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}

	count, _ := strconv.Atoi(ctx.Query("count"))
	if count == 0 {
		count = 10
	}

	messages, total, err := t.tgMessageService.GetList(
		ctx,
		ctx.Query("chat_id"),
		page,
		count,
	)
	if err != nil {
		logger.OPError(op, err)
		response.Fail(ctx, err)
		return
	}

	response.Pagination(ctx, page, count, total, messages)
}
