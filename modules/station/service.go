package station

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/common/client"
)

// Service adalah "kontrak" (interface) yang menentukan fungsi apa saja
// yang harus dimiliki oleh service station.
// Di sini cuma ada 1 fungsi: GetAllStation.
type Service interface {
	GetAllStation() (resp []StationOut, err error)
}

// service adalah implementasi dari Service.
// Struct ini punya field "client" untuk melakukan HTTP request.
type service struct {
	client *http.Client
}

// NewService membuat object service baru.
// Di sini kita juga set timeout untuk HTTP client supaya request tidak menggantung terlalu lama.
func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetAllStation memanggil API MRT (https://jakartamrt.co.id/id/val/stasiuns)
// untuk mengambil daftar stasiun.
// 1. Panggil API → dapat response dalam bentuk byte.
// 2. Ubah byte jadi slice of StationIn (pakai json.Unmarshal).
// 3. Konversi StationIn → StationOut.
// 4. Kembalikan hasilnya ke pemanggil.
func (s *service) GetAllStation() (resp []StationOut, err error) {

	// URL API eksternal MRT
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	// Lakukan HTTP GET ke API
	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	// Simpan hasil parsing dari API ke slice of StationIn
	var stations []StationIn
	err = json.Unmarshal(byteResponse, &stations)

	// Ubah setiap StationIn menjadi StationOut
	for _, item := range stations {
		resp = append(resp, StationOut(item))
	}

	return
}
