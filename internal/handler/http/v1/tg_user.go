package v1

import (
	"fmt"
	"strconv"

	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/gin-gonic/gin"
)

type tgUser struct {
	l             logger.Logger
	tgUserService *service.TGUser
}

func newTGUser(l logger.Logger, router *gin.RouterGroup, tgUserService *service.TGUser) *tgUser {
	t := &tgUser{
		l:             l,
		tgUserService: tgUserService,
	}

	router.GET("/tg-users", t.all)

	return t
}

// @Summary Show users
// @Description Show all users list
// @Tags tg-users
// @Accept json
// @Produce json
// @Param username query string false "Username"
// @Param chat_id query string false "ChatID"
// @Param page query int false "Page"
// @Param count query int false "Count"
// @Success 200 {object} successResponse{data=[]entity.TGUser}
// @Failure 500 {object} errorResponse
// @Router /tg-users [get]
func (t *tgUser) all(ctx *gin.Context) {
	const op = "./internal/handler/http/v1/tg_user::all"

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		count = 0
	}

	users, err := t.tgUserService.All(ctx, ctx.Query("username"), ctx.Query("chat_id"), page, count)
	if err != nil {
		t.l.Error(fmt.Errorf("%s: %w", op, err))
		errorRsp(ctx, err)
		return
	}

	successRsp(ctx, users)
	return
}
