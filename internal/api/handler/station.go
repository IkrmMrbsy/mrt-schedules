package handler

import (
	"net/http"

	"github.com/IkrmMrbsy/mrt-schedules/internal/api/service/station"
	"github.com/IkrmMrbsy/mrt-schedules/pkg/response"
	"github.com/gin-gonic/gin"
)

// Initiate digunakan untuk mendaftarkan semua route (endpoint) terkait station.
// - Pertama buat service station (pakai NewService).
// - Lalu daftarkan route /stations GET yang akan memanggil fungsi GetAllStation.
func Initiate(router *gin.RouterGroup, stationService station.Service) {

	// Buat group route "/stations"
	station := router.Group("/stations")

	// GET /stations
	station.GET("/", func(ctx *gin.Context) {
		GetAllStation(ctx, stationService)
	})

	station.GET("/:id", func(ctx *gin.Context) {
		CheckScheduleByStation(ctx, stationService)
	})

	station.GET("/fare", func(ctx *gin.Context) {
		GetFareAndDuration(ctx, stationService)
	})
}

// GetAllStation adalah handler untuk route GET /stations.
// Handler = fungsi yang akan dijalankan ketika endpoint dipanggil.
// 1. Panggil service.GetAllStation() → ambil data stasiun dari API MRT.
// 2. Kalau error, balikin response 400 (Bad Request).
// 3. Kalau sukses, balikin response 200 (OK) beserta data stasiun.
func GetAllStation(ctx *gin.Context, service station.Service) {
	datas, err := service.GetAllStation()
	if err != nil {
		// Jika error, kembalikan HTTP 400
		ctx.JSON(http.StatusBadRequest,
			response.APIResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}

	// Jika sukses, kembalikan HTTP 200 dengan data stasiun
	ctx.JSON(http.StatusOK,
		response.APIResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    datas,
		},
	)
}

func CheckScheduleByStation(ctx *gin.Context, service station.Service) {
	id := ctx.Param("id")

	datas, err := service.CheckScheduleByStation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			response.APIResponse{
				Code:    http.StatusBadGateway,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}

	ctx.JSON(http.StatusOK,
		response.APIResponse{
			Code:    http.StatusOK,
			Message: "Success",
			Data:    datas,
		},
	)
}

func GetFareAndDuration(ctx *gin.Context, service station.Service) {
	fromId := ctx.Query("from")
	toId := ctx.Query("to")

	data, err := service.GetFareAndDuration(fromId, toId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest,
			response.APIResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}

	ctx.JSON(http.StatusOK,
		response.APIResponse{
			Code:    http.StatusOK,
			Message: "success",
			Data:    data,
		},
	)
}
