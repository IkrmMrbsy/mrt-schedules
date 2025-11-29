package station

import (
	"errors"
	"strings"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/internal/api/service/station"
	"github.com/IkrmMrbsy/mrt-schedules/pkg/utils"
)

type Usecase interface {
	GetAllStation(name string) ([]StationOut, error)
	CheckScheduleByStation(id string) ([]ScheduleOut, error)
	GetFareAndDuration(fromId, toId string) (FareOut, error)
	GetNextTrainByStation(id, destination string) (*NextTrainOut, error)
	GetStationDetails(id string) (*DetailStationOut, error)
}

type usecase struct {
	service station.Service
}

func NewUsecase(service station.Service) Usecase {
	return &usecase{service: service}
}

func (u *usecase) GetAllStation(name string) ([]StationOut, error) {
	stations, err := u.service.FetchStations()
	if err != nil {
		return nil, err
	}

	if name != "" {
		stations = utils.Filter(stations, func(s station.StationIn) bool {
			return strings.Contains(strings.ToLower(s.NamaStasiun), strings.ToLower(name))
		})
	}

	var resp []StationOut
	for _, item := range stations {
		resp = append(resp, StationOut{
			Id:   item.ID,
			Nama: item.NamaStasiun,
		})
	}

	return resp, nil
}

func (u *usecase) CheckScheduleByStation(id string) ([]ScheduleOut, error) {
	schedules, err := u.service.FetchSchedules()
	if err != nil {
		return nil, err
	}

	var scheduleSelected station.ScheduleIn
	for _, item := range schedules {
		if item.IDStasiun == id {
			scheduleSelected = item
			break
		}
	}
	if scheduleSelected.IDStasiun == "" {
		return nil, errors.New("station not found")
	}

	return ConvertDataToResponse(scheduleSelected)
}

func (u *usecase) GetFareAndDuration(fromId, toId string) (FareOut, error) {
	stations, err := u.service.FetchFares()
	if err != nil {
		return FareOut{}, err
	}

	var fromName, toName, fare, duration string
	for _, st := range stations {
		if st.ID == fromId {
			fromName = st.Nama
			for _, e := range st.Estimasi {
				if e.IDStasiunTujuan == toId {
					fare = e.Tarif
					duration = e.Waktu + " menit"
					break
				}
			}
		}
		if st.ID == toId {
			toName = st.Nama
		}
		if fromName != "" && toName != "" && fare != "" {
			break
		}
	}

	if fromName == "" || toName == "" {
		return FareOut{}, errors.New("station not found")
	}
	if fare == "" {
		return FareOut{}, errors.New("fare/estimasi not found between stations")
	}

	return FareOut{
		Dari:   fromName,
		Ke:     toName,
		Tarif:  fare,
		Durasi: duration,
	}, nil
}

func (u *usecase) GetNextTrainByStation(id, destination string) (*NextTrainOut, error) {
	schedules, err := u.service.FetchSchedules()
	if err != nil {
		return nil, err
	}

	var scheduleSelected station.ScheduleIn
	for _, item := range schedules {
		if item.IDStasiun == id {
			scheduleSelected = item
			break
		}
	}
	if scheduleSelected.IDStasiun == "" {
		return nil, errors.New("station not found")
	}

	isWeekend := time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday
	var times []time.Time

	if destination == "LB" {
		if isWeekend {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.JadwalLebakBulusLibur)
		} else {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.JadwalLebakBulusBiasa)
		}
	} else if destination == "HI" {
		if isWeekend {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.JadwalBundaranHILibur)
		} else {
			times, err = ConvertScheduleToTimeFormat(scheduleSelected.JadwalBundaranHIBiasa)
		}
	} else {
		return nil, errors.New("invalid destination, use 'LB' or 'HI'")
	}
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var nextTrains []TrainSchedule
	for _, t := range times {
		if t.After(now) {
			nextTrains = append(nextTrains, TrainSchedule{WaktuKeberangkatan: t.Format("15:04")})
			if len(nextTrains) == 3 {
				break
			}
		}
	}

	if len(nextTrains) == 0 {
		return nil, errors.New("no next train available today")
	}

	return &NextTrainOut{
		IdKereta:         id,
		Stasiun:          scheduleSelected.NamaStasiun,
		Tujuan:           DestinationMap[destination],
		KeretaBerikutnya: nextTrains,
	}, nil
}

func (u *usecase) GetStationDetails(id string) (*DetailStationOut, error) {
	stations, err := u.service.FetchStations()
	if err != nil {
		return nil, err
	}

	var stationData *station.StationIn
	for _, st := range stations {
		if st.ID == id {
			stationData = &st
			break
		}
	}

	if stationData == nil {
		return nil, errors.New("station not found")
	}

	antarmodaParsed := ParseAntarmoda(stationData.Antarmoda)

	komersial := GroupRetailAndFacilities(stationData.Retails, stationData.Fasilitas)

	resp := &DetailStationOut{
		ID:          stationData.ID,
		NamaStasiun: stationData.NamaStasiun,
		Gambar: GambarOut{
			Banner:        stationData.Banner,
			PetaLokalitas: stationData.PetaLokalitas,
		},
		TransportasiLanjutan: antarmodaParsed,
		FasilitasKomersial:   komersial,
	}

	return resp, nil
}
