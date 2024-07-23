package v1

import (
	"fmt"
	"time"

	_ "github.com/eerzho/ten_tarot/docs"
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
func NewHandler(router *gin.Engine, tgUserService tgUserService, tgMessageService tgMessageService) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	router.GET("/swagger/*any", swaggerHandler)

	mv := newMiddleware()
	router.Use(mv.log)

	v1 := router.Group("/api/v1")
	{
		newTGUser(v1, tgUserService)
		newTGMessage(v1, tgMessageService)
	}
}

type middleware struct{}

func newMiddleware() *middleware {
	return &middleware{}
}

func (m *middleware) log(ctx *gin.Context) {
	id := uuid.New().String()
	ctx.Set(RID, id)

	logger.Info(fmt.Sprintf("start: %s", id))

	start := time.Now()
	ctx.Next()
	duration := time.Since(start)

	logger.Info(fmt.Sprintf("end: %s - %.4f sec.", id, duration.Seconds()))
}
