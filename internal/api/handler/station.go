package handler

import (
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

	station.GET("/:id/next-train", func(ctx *gin.Context) {
		GetNextTrainByStation(ctx, stationService)
	})
}

// GetAllStation adalah handler untuk route GET /stations.
// Handler = fungsi yang akan dijalankan ketika endpoint dipanggil.
// 1. Panggil service.GetAllStation() → ambil data stasiun dari API MRT.
// 2. Kalau error, balikin response 400 (Bad Request).
// 3. Kalau sukses, balikin response 200 (OK) beserta data stasiun.
func GetAllStation(ctx *gin.Context, service station.Service) {
	resp, err := service.GetAllStation()
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// Jika sukses, kembalikan HTTP 200 dengan data stasiun
	response.Success(ctx, resp)
}

func CheckScheduleByStation(ctx *gin.Context, service station.Service) {
	id := ctx.Param("id")

	resp, err := service.CheckScheduleByStation(id)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func GetFareAndDuration(ctx *gin.Context, service station.Service) {
	fromId := ctx.Query("from")
	toId := ctx.Query("to")

	resp, err := service.GetFareAndDuration(fromId, toId)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func GetNextTrainByStation(ctx *gin.Context, service station.Service) {
	id := ctx.Param("id")
	destination := ctx.Query("destination")

	resp, err := service.GetNextTrainByStation(id, destination)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}
