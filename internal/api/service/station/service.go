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
	GetFareAndDuration(fromId, toId string) (resp FareOut, err error)
	GetNextTrainByStation(id, destination string) (resp *NextTrainOut, err error)
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

func (s *service) GetFareAndDuration(fromId, toId string) (resp FareOut, err error) {
	byteResp, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return resp, err
	}

	var stations []FareIn
	if err := json.Unmarshal(byteResp, &stations); err != nil {
		return resp, err
	}

	var fromName, toName string
	var fare, duration string

	for _, st := range stations {
		if st.Id == fromId {
			fromName = st.Name
			for _, e := range st.Estimasi {
				if e.StationNid == toId {
					fare = e.Tarif
					duration = e.Waktu + " menit"
					break
				}
			}
		}
		if st.Id == toId {
			toName = st.Name
		}
		if fromName != "" && toName != "" && fare != "" {
			break
		}
	}

	if fare == "" {
		for _, st := range stations {
			if st.Id == toId {
				for _, e := range st.Estimasi {
					if e.StationNid == fromId {
						fare = e.Tarif
						duration = e.Waktu
						break
					}
				}
			}
			if fare != "" {
				break
			}
		}
	}

	if fromName == "" || toName == "" {
		return resp, errors.New("station not found")
	}
	if fare == "" {
		return resp, errors.New("fare/estimasi not found between stations")
	}

	resp.From = fromName
	resp.To = toName
	resp.Fare = fare
	resp.Duration = duration

	return resp, nil
}

func (s *service) GetNextTrainByStation(id, destination string) (resp *NextTrainOut, err error) {
	byteResp, err := client.DoRequest(s.client, s.apiURL)
	if err != nil {
		return resp, nil
	}

	var schedules []ScheduleIn
	if err := json.Unmarshal(byteResp, &schedules); err != nil {
		return resp, nil
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

	isWeekend := time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday

	var times []time.Time

	if destination == "LB" {
		if isWeekend {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.ScheduleLebakBulusLibur)
		} else {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.ScheduleLebakBulus)
		}
	} else if destination == "HI" {
		if isWeekend {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.ScheduleBundaranHILibur)
		} else {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.ScheduleBundaranHI)
		}
	} else {
		return nil, errors.New("invalid destination, use 'LB' or 'HI'")
	}

	if err != nil {
		return nil, err
	}

	now := time.Now()
	var next time.Time
	for _, t := range times {
		if t.After(now) && (next.IsZero() || t.Before(next)) {
			next = t
		}
	}

	if next.IsZero() {
		return nil, errors.New("no next train available today")
	}

	resp = &NextTrainOut{
		TrainId:     id,
		Destination: scheduleSelected.StationName,
		Departure:   next.Format("15:04"),
	}

	return resp, nil
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
	var schedules = strings.Split(schedule, ",")
	now := time.Now()
	today := now.Format("2006-01-02")

	loc := time.Now().Location()

	for _, item := range schedules {
		trimedTime := strings.TrimSpace(item)
		if trimedTime == "" {
			continue
		}

		fullTimeStr := today + " " + trimedTime
		parsed, err := time.ParseInLocation("2006-01-02 15:04", fullTimeStr, loc)
		if err != nil {
			err = errors.New("invalid time format " + trimedTime)
			return nil, err
		}
		resp = append(resp, parsed)
	}

	return
}
