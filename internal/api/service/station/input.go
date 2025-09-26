package station

// StationIn dipakai untuk menampung data dari API eksternal (API MRT).
// Field "Id" dan "Name" akan otomatis diisi dari JSON yang dikirim API.
// Tag `json:"nid"` artinya Id diisi dari field "nid" pada JSON.
// Tag `json:"title"` artinya Name diisi dari field "title" pada JSON.
type StationIn struct {
	Id   string `json:"nid"`
	Name string `json:"title"`
}

type ScheduleIn struct {
	StationId          string `json:"nid"`
	StationName        string `json:"title"`
	ScheduleBundaranHI string `json:"jadwal_hi_biasa"`
	ScheduleLebakBulus string `json:"jadwal_lb_biasa"`
}
