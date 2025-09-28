package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, APIError{
		Code:    status,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(ctx *gin.Context, message string) {
	Error(ctx, http.StatusBadRequest, message)
}

func NotFound(ctx *gin.Context, message string) {
	Error(ctx, http.StatusNotFound, message)
}
