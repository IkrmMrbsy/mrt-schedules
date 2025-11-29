package station

import (
	"errors"
	"strings"
	"time"

	"github.com/IkrmMrbsy/mrt-schedules/internal/api/service/station"
)

func ConvertDataToResponse(schedule station.ScheduleIn) (resp []ScheduleOut, err error) {
	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(schedule.JadwalLebakBulusBiasa)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(schedule.JadwalBundaranHIBiasa)
	if err != nil {
		return
	}

	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			resp = append(resp, ScheduleOut{
				NamaStasiun: schedule.NamaStasiun,
				Waktu:       item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			resp = append(resp, ScheduleOut{
				NamaStasiun: schedule.NamaStasiun,
				Waktu:       item.Format("15:04"),
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
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", fullTimeStr, loc)
		if err != nil {
			err = errors.New("invalid time format " + trimedTime)
			return nil, err
		}
		resp = append(resp, parsed)
	}

	return
}

func ParseAntarmoda(antarmodaStr string) []AntarmodaOut {
	if antarmodaStr == "" {
		return nil
	}

	// Pisahkan berdasarkan baris (Metromini\r\nB85, S72\r\n\r\nKWK\r\nS03, ...)
	parts := strings.Split(antarmodaStr, "\r\n\r\n")
	var results []AntarmodaOut

	for _, part := range parts {
		lines := strings.Split(strings.TrimSpace(part), "\r\n")
		if len(lines) < 2 {
			continue
		}

		jenis := strings.TrimSpace(lines[0])
		ruteStr := strings.TrimSpace(lines[1])
		rute := strings.Split(ruteStr, ", ")

		// Membersihkan rute dari spasi dan karakter sisa
		var cleanRute []string
		for _, r := range rute {
			if cleaned := strings.TrimSpace(r); cleaned != "" {
				cleanRute = append(cleanRute, cleaned)
			}
		}

		if len(cleanRute) > 0 {
			results = append(results, AntarmodaOut{
				Jenis: jenis,
				Rute:  cleanRute,
			})
		}
	}
	return results
}

// GroupRetailAndFacilities menggabungkan Retail dan Fasilitas ke dalam map berdasarkan jenis/tipe.
func GroupRetailAndFacilities(retails []station.RetailIn, fasilitas []station.FasilitasIn) map[string][]FasilitasOut {
	grouped := make(map[string][]FasilitasOut)

	// Proses Retail
	for _, r := range retails {
		// Tentukan tipe, jika null/kosong gunakan "Lain-lain" atau kategorisasi sendiri.
		jenis := strings.TrimSpace(r.JenisRetail)
		if jenis == "" {
			jenis = "Lain-lain"
		}

		// Normalisasi dan kapitalisasi jenis
		jenis = strings.Title(strings.ToLower(jenis))

		item := FasilitasOut{
			ID:    r.ID,
			Nama:  r.Judul,
			Cover: r.Cover,
			Tipe:  jenis,
		}
		grouped[jenis] = append(grouped[jenis], item)
	}

	// Proses Fasilitas
	for _, f := range fasilitas {
		jenis := strings.TrimSpace(f.JenisFasilitas)
		if jenis == "" {
			jenis = "Lain-lain"
		}

		jenis = strings.Title(strings.ToLower(jenis))

		item := FasilitasOut{
			ID:    f.ID,
			Nama:  f.Judul,
			Cover: f.Cover,
			Tipe:  jenis,
		}
		grouped[jenis] = append(grouped[jenis], item)
	}

	return grouped
}
