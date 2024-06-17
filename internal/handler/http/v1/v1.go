package v1

import (
	_ "github.com/eerzho/event_manager/docs"
	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHandler -.
// Swagger spec:
// @Title       Event manager tg bot api
// @Version     1.0
// @BasePath    /api/v1
func NewHandler(l logger.Logger, router *gin.Engine, tgUserService *service.TGUser, tgMessageService *service.TGMessage) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	router.GET("/swagger/*any", swaggerHandler)

	v1 := router.Group("/api/v1")
	{
		newTGUser(l, v1, tgUserService)
		newTGMessage(l, v1, tgMessageService)
	}
}
