package main

import (
	"github.com/IkrmMrbsy/mrt-schedules/internal/api/handler"
	"github.com/IkrmMrbsy/mrt-schedules/internal/api/service/station"
	stationUsecase "github.com/IkrmMrbsy/mrt-schedules/internal/api/usecase/station"
	"github.com/IkrmMrbsy/mrt-schedules/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	stationService := station.NewService(cfg.HttpTimeout, cfg.MRTApiURL)
	stationUsecase := stationUsecase.NewUsecase(stationService)

	// Jalankan fungsi InitiateRoutes untuk memulai server
	InitiateRoutes(stationUsecase, cfg.ServerPort)
}

// InitiateRoutes bertugas untuk:
// 1. Membuat router baru (pakai Gin).
// 2. Membuat group endpoint dengan prefix "/v1/api".
// 3. Daftarkan semua route dari module station.
// 4. Menjalankan server di port 8080.
func InitiateRoutes(stationUsecase stationUsecase.Usecase, port string) {
	var (
		router = gin.Default()           // router utama (sudah ada logger + recovery bawaan)
		api    = router.Group("/v1/api") // prefix semua route diawali /v1/api
	)

	// Daftarkan semua endpoint station ke dalam group /v1/api
	handler.Initiate(api, stationUsecase)

	// Jalankan server di port 8080
	router.Run(":" + port)
}
