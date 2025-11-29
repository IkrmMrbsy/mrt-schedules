package handler

import (
	"github.com/IkrmMrbsy/mrt-schedules/internal/api/usecase/station"
	"github.com/IkrmMrbsy/mrt-schedules/pkg/response"
	"github.com/gin-gonic/gin"
)

// Initiate digunakan untuk mendaftarkan semua route (endpoint) terkait station.
// - Pertama buat service station (pakai NewService).
// - Lalu daftarkan route /stations GET yang akan memanggil fungsi GetAllStation.
func Initiate(router *gin.RouterGroup, usecase station.Usecase) {

	// Buat group route "/stations"
	station := router.Group("/stations")

	// GET /stations
	station.GET("/", func(ctx *gin.Context) {
		GetAllStation(ctx, usecase)
	})

	station.GET("/:id", func(ctx *gin.Context) {
		CheckScheduleByStation(ctx, usecase)
	})

	station.GET("/fare", func(ctx *gin.Context) {
		GetFareAndDuration(ctx, usecase)
	})

	station.GET("/:id/next-train", func(ctx *gin.Context) {
		GetNextTrainByStation(ctx, usecase)
	})

	station.GET("/:id/details", func(ctx *gin.Context) {
		GetStationDetails(ctx, usecase)
	})
}

// GetAllStation adalah handler untuk route GET /stations.
// Handler = fungsi yang akan dijalankan ketika endpoint dipanggil.
// 1. Panggil service.GetAllStation() → ambil data stasiun dari API MRT.
// 2. Kalau error, balikin response 400 (Bad Request).
// 3. Kalau sukses, balikin response 200 (OK) beserta data stasiun.
func GetAllStation(ctx *gin.Context, usecase station.Usecase) {
	name := ctx.Query("name")

	resp, err := usecase.GetAllStation(name)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// Jika sukses, kembalikan HTTP 200 dengan data stasiun
	response.Success(ctx, resp)
}

func CheckScheduleByStation(ctx *gin.Context, usecase station.Usecase) {
	id := ctx.Param("id")

	resp, err := usecase.CheckScheduleByStation(id)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func GetFareAndDuration(ctx *gin.Context, usecase station.Usecase) {
	fromId := ctx.Query("from")
	toId := ctx.Query("to")

	resp, err := usecase.GetFareAndDuration(fromId, toId)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func GetNextTrainByStation(ctx *gin.Context, usecase station.Usecase) {
	id := ctx.Param("id")
	destination := ctx.Query("destination")

	resp, err := usecase.GetNextTrainByStation(id, destination)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func GetStationDetails(ctx *gin.Context, usecase station.Usecase) {
	id := ctx.Param("id")

	resp, err := usecase.GetStationDetails(id)
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, resp)
}
