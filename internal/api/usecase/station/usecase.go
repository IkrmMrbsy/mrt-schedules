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
			return strings.Contains(strings.ToLower(s.Name), strings.ToLower(name))
		})
	}

	var resp []StationOut
	for _, item := range stations {
		resp = append(resp, StationOut(item))
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
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}
	if scheduleSelected.StationId == "" {
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

	if fromName == "" || toName == "" {
		return FareOut{}, errors.New("station not found")
	}
	if fare == "" {
		return FareOut{}, errors.New("fare/estimasi not found between stations")
	}

	return FareOut{
		From:     fromName,
		To:       toName,
		Fare:     fare,
		Duration: duration,
	}, nil
}

func (u *usecase) GetNextTrainByStation(id, destination string) (*NextTrainOut, error) {
	schedules, err := u.service.FetchSchedules()
	if err != nil {
		return nil, err
	}

	var scheduleSelected station.ScheduleIn
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
	var nextTrains []TrainSchedule
	for _, t := range times {
		if t.After(now) {
			nextTrains = append(nextTrains, TrainSchedule{Departure: t.Format("15:04")})
			if len(nextTrains) == 3 {
				break
			}
		}
	}

	if len(nextTrains) == 0 {
		return nil, errors.New("no next train available today")
	}

	return &NextTrainOut{
		TrainId:     id,
		Station:     scheduleSelected.StationName,
		Destination: DestinationMap[destination],
		NextTrains:  nextTrains,
	}, nil
}
