package main

import (
	"github.com/IkrmMrbsy/mrt-schedules/modules/station"
	"github.com/gin-gonic/gin"
)

func main() {
	// Jalankan fungsi InitiateRoutes untuk memulai server
	InitiateRoutes()
}

// InitiateRoutes bertugas untuk:
// 1. Membuat router baru (pakai Gin).
// 2. Membuat group endpoint dengan prefix "/v1/api".
// 3. Daftarkan semua route dari module station.
// 4. Menjalankan server di port 8080.
func InitiateRoutes() {
	var (
		router = gin.Default()           // router utama (sudah ada logger + recovery bawaan)
		api    = router.Group("/v1/api") // prefix semua route diawali /v1/api
	)

	// Daftarkan semua endpoint station ke dalam group /v1/api
	station.Initiate(api)

	// Jalankan server di port 8080
	router.Run(":8080")
}
