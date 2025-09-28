package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse adalah struct standar untuk balikan (response) API.
// Semua endpoint API akan memakai format ini supaya konsisten.
// - Code    → angka status HTTP (contoh: 200, 400, 500).
// - Message → pesan singkat tentang hasil request.
// - Data    → isi data utama (bisa apa saja: list, object, atau nil).
type APISuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, APISuccess{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}
