package station

import (
	"errors"
	"strings"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/internal/api/service/station"
)

func ConvertDataToResponse(schedule station.ScheduleIn) (resp []ScheduleOut, err error) {
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
