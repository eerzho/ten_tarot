package v1

import (
	"fmt"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/handler/http/v1/response"
	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/logger"
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

	router.GET("/tg-users", t.list)

	return t
}

// list -.
// @Summary Show users
// @Description Show users list
// @Tags tg-users
// @Accept json
// @Produce json
// @Param username query string false "Username"
// @Param chat_id query string false "ChatID"
// @Param page query int false "Page"
// @Param count query int false "Count"
// @Success 200 {object} response.pagination{data=[]entity.TGUser}
// @Router /tg-users [get]
func (t *tgUser) list(ctx *gin.Context) {
	const op = "./internal/handler/http/v1/tg_user::list"

	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}

	count, _ := strconv.Atoi(ctx.Query("count"))
	if count == 0 {
		count = 10
	}

	users, total, err := t.tgUserService.List(ctx, ctx.Query("username"), ctx.Query("chat_id"), page, count)
	if err != nil {
		t.l.Error(fmt.Sprintf("%s - %s", op, err.Error()))
		response.Fail(ctx, err)
		return
	}

	response.Pagination(ctx, page, count, total, users)
}
