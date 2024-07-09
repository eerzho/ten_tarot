package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eerzho/ten_tarot/internal/entity"
	"github.com/eerzho/ten_tarot/internal/failure"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	fail struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message"`
	}

	success struct {
		Data interface{} `json:"data"`
	}

	pagination struct {
		Data       interface{}        `json:"data"`
		Pagination *entity.Pagination `json:"pagination,omitempty"`
	}
)

func Fail(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		code = http.StatusBadRequest
		data := make(map[string]string)
		for _, v := range ve {
			data[v.Field()] = fmt.Sprintf("failed on the '%s' tag", v.Tag())
		}
		ctx.AbortWithStatusJSON(code, fail{Data: data, Message: "validation failed"})
		return
	} else if errors.Is(err, failure.ErrNotFound) {
		code = http.StatusNotFound
	}

	ctx.AbortWithStatusJSON(code, fail{Message: err.Error()})
}

func Success(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	ctx.JSON(http.StatusOK, success{data})
}

func Pagination(ctx *gin.Context, currentPage, countPerPage, total int, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	ctx.JSON(http.StatusOK, pagination{
		Data:       data,
		Pagination: &entity.Pagination{CurrentPage: currentPage, CountPerPage: countPerPage, Total: total},
	})
}
