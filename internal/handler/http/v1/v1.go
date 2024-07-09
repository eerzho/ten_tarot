package v1

import (
	"fmt"
	"time"

	_ "github.com/eerzho/ten_tarot/docs"
	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const RID = "X-Request-ID"

// NewHandler -.
// Swagger spec:
// @Title       Ten tarot tg bot api
// @Version     1.0
// @BasePath    /api/v1
func NewHandler(l logger.Logger, router *gin.Engine, tgUserService *service.TGUser, tgMessageService *service.TGMessage) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	router.GET("/swagger/*any", swaggerHandler)

	mv := newMiddleware(l)
	router.Use(mv.log)

	v1 := router.Group("/api/v1")
	{
		newTGUser(l, v1, tgUserService)
		newTGMessage(l, v1, tgMessageService)
	}
}

type middleware struct {
	l logger.Logger
}

func newMiddleware(l logger.Logger) *middleware {
	return &middleware{
		l: l,
	}
}

func (m *middleware) log(ctx *gin.Context) {
	id := uuid.New().String()
	ctx.Set(RID, id)

	m.l.Info(fmt.Sprintf("start: %s", id))

	start := time.Now()
	ctx.Next()
	duration := time.Since(start)

	m.l.Info(fmt.Sprintf("end: %s - %s", id, duration.String()))
}
