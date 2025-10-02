package station

// StationOut dipakai untuk mengirim data balik (response) ke API kita sendiri.
// Bedanya dengan StationIn ada di nama field JSON yang kita tentukan.
// Tag `json:"id"` artinya Id dikirim dengan nama "id".
// Tag `json:"name"` artinya Name dikirim dengan nama "name".
type StationOut struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ScheduleOut struct {
	StationName string `json:"title"`
	Time        string `json:"time"`
}

type FareOut struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Fare     string `json:"fare"`
	Duration string `json:"duration"`
}

type NextTrainOut struct {
	TrainId     string          `json:"train_id"`
	Station     string          `json:"station"`
	Destination string          `json:"destination"`
	NextTrains  []TrainSchedule `json:"next_trains"`
}

type TrainSchedule struct {
	Departure string `json:"deperture_time"`
}
