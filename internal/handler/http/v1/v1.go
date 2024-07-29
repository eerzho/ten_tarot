package v1

import (
	_ "github.com/eerzho/ten_tarot/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewHandler -.
// Swagger spec:
// @Title       Ten tarot tg bot api
// @Version     1.0
// @BasePath    /api/v1
func NewHandler(
	router *gin.Engine,
	tgUserService tgUserService,
	tgMessageService tgMessageService,
) {
	swaggerHandler := ginSwagger.DisablingWrapHandler(
		swaggerFiles.Handler,
		"DISABLE_SWAGGER_HTTP_HANDLER",
	)
	router.GET("/swagger/*any", swaggerHandler)

	mv := newMiddleware()

	router.Use(mv.setRIDAndLogDuration)

	v1 := router.Group("/api/v1")

	newTGUser(v1, tgUserService)
	newTGMessage(v1, tgMessageService)
}
