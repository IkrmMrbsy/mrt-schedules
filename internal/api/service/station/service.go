package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/pkg/client"
)

// Service adalah "kontrak" (interface) yang menentukan fungsi apa saja
// yang harus dimiliki oleh service station.
// Di sini cuma ada 1 fungsi: GetAllStation.
type Service interface {
	GetAllStation() (resp []StationOut, err error)
	CheckScheduleByStation(id string) (resp []ScheduleOut, err error)
}

// service adalah implementasi dari Service.
// Struct ini punya field "client" untuk melakukan HTTP request.
type service struct {
	client *http.Client
	apiURL string
}

// NewService membuat object service baru.
// Di sini kita juga set timeout untuk HTTP client supaya request tidak menggantung terlalu lama.
func NewService(apiURL string) Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
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
func (s *service) GetAllStation() (resp []StationOut, err error) {
	// Lakukan HTTP GET ke API
	byteResponse, err := client.DoRequest(s.client, s.apiURL)
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

func (s *service) CheckScheduleByStation(id string) (resp []ScheduleOut, err error) {
	byteResponse, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return
	}

	var schedules []ScheduleIn
	if err := json.Unmarshal(byteResponse, &schedules); err != nil {
		return nil, err
	}

	var scheduleSelected ScheduleIn
	for _, item := range schedules {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		return nil, errors.New("station not found")
	}

	resp, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}

	return
}

func ConvertDataToResponse(schedule ScheduleIn) (resp []ScheduleOut, err error) {
	// var (
	// 	LebakBulusTripName = "Stasiun Lebak Bulus Grab"
	// 	BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	// )

	// scheduleLebakBulus := schedule.ScheduleLebakBulus
	// scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleBundaranHI)
	if err != nil {
		return
	}

	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			resp = append(resp, ScheduleOut{
				StationName: schedule.StationName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			resp = append(resp, ScheduleOut{
				StationName: schedule.StationName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (resp []time.Time, err error) {
	var (
		ParsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimedTime := strings.TrimSpace(item)
		if trimedTime == "" {
			continue
		}

		ParsedTime, err = time.Parse("15:04", trimedTime)
		if err != nil {
			err = errors.New("invalid time format " + trimedTime)
			return
		}
		resp = append(resp, ParsedTime)
	}

	return
}
