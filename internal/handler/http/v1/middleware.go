package v1

import (
	"time"

	"github.com/eerzho/ten_tarot/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RID = "X-Request-ID"

type middleware struct{}

func newMiddleware() *middleware {
	return &middleware{}
}

func (m *middleware) setRIDAndLogDuration(ctx *gin.Context) {
	const op = "handler.http.v1.middleware.setRIDAndLogDuration"

	id := uuid.New().String()
	ctx.Set(RID, id)

	logger.Info(op, logger.Any("id", id))

	start := time.Now()
	ctx.Next()
	duration := time.Since(start)

	logger.Info(
		op,
		logger.Any("id", id),
		logger.Any("duration in seconds", duration.Seconds()),
	)
}
