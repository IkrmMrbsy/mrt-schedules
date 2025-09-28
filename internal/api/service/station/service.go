package station

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/pkg/client"
)

// Service adalah "kontrak" (interface) yang menentukan fungsi apa saja
// yang harus dimiliki oleh service station.
// Di sini cuma ada 1 fungsi: GetAllStation.
type Service interface {
	FetchStations() ([]StationIn, error)
	FetchSchedules() ([]ScheduleIn, error)
	FetchFares() ([]FareIn, error)
}

// service adalah implementasi dari Service.
// Struct ini punya field "client" untuk melakukan HTTP request.
type service struct {
	client *http.Client
	apiURL string
}

// NewService membuat object service baru.
// Di sini kita juga set timeout untuk HTTP client supaya request tidak menggantung terlalu lama.
func NewService(timeout time.Duration, apiURL string) Service {
	return &service{
		client: &http.Client{
			Timeout: timeout,
		},
		apiURL: apiURL,
	}
}

// GetAllStation memanggil API MRT (https://jakartamrt.co.id/id/val/stasiuns)
// untuk mengambil daftar stasiun.
// 1. Panggil API → dapat response dalam bentuk byte.
// 2. Ubah byte jadi slice of StationIn (pakai json.Unmarshal).
// 3. Konversi StationIn → StationOut.
// 4. Kembalikan hasilnya ke pemanggil.
func (s *service) FetchStations() ([]StationIn, error) {
	// Lakukan HTTP GET ke API
	byteResponse, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return nil, err
	}

	// Simpan hasil parsing dari API ke slice of StationIn
	var stations []StationIn
	err = json.Unmarshal(byteResponse, &stations)
	if err != nil {
		return nil, err
	}

	return stations, nil
}

func (s *service) FetchSchedules() ([]ScheduleIn, error) {
	byteResponse, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return nil, err
	}

	var schedules []ScheduleIn
	if err := json.Unmarshal(byteResponse, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *service) FetchFares() ([]FareIn, error) {
	byteResp, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return nil, err
	}

	var stations []FareIn
	if err := json.Unmarshal(byteResp, &stations); err != nil {
		return nil, err
	}

	return stations, nil
}
