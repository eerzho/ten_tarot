package v1

import (
	"errors"
	"net/http"

	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func errorRsp(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	if errors.Is(err, failure.ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, failure.ErrValidation) {
		code = http.StatusBadRequest
	}

	ctx.AbortWithStatusJSON(code, errorResponse{err.Error()})
}

type successResponse struct {
	Data interface{} `json:"data"`
}

func successRsp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, successResponse{data})
}
