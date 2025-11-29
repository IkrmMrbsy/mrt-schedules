package station

// StationOut (Output Data Ringkas Stasiun)
type StationOut struct {
	Id   string `json:"id"`
	Nama string `json:"nama"`
}

// ScheduleOut (Output Jadwal Per Keberangkatan)
type ScheduleOut struct {
	NamaStasiun string `json:"nama_stasiun"`
	Waktu       string `json:"waktu"`
}

// FareOut (Output Tarif dan Durasi)
type FareOut struct {
	Dari   string `json:"dari"`
	Ke     string `json:"ke"`
	Tarif  string `json:"tarif"`
	Durasi string `json:"durasi"` // Contoh: "10 menit"
}

// TrainSchedule (Sub-struct untuk Waktu Keberangkatan)
type TrainSchedule struct {
	WaktuKeberangkatan string `json:"waktu_keberangkatan"`
}

// NextTrainOut (Output Kereta Berikutnya)
type NextTrainOut struct {
	IdKereta         string          `json:"id_kereta"`
	Stasiun          string          `json:"stasiun"`
	Tujuan           string          `json:"tujuan"`
	KeretaBerikutnya []TrainSchedule `json:"kereta_berikutnya"`
}

// DetailStationOut (Output Detail Lengkap Stasiun)
type DetailStationOut struct {
	ID                   string                    `json:"id"`
	NamaStasiun          string                    `json:"nama_stasiun"`
	Gambar               GambarOut                 `json:"gambar"` // Ubah nama field JSON jadi lebih ringkas
	TransportasiLanjutan []AntarmodaOut            `json:"transportasi_lanjutan"`
	FasilitasKomersial   map[string][]FasilitasOut `json:"fasilitas_komersial"`
}

// GambarOut (Sub-struct untuk Informasi Visual)
type GambarOut struct {
	Banner        string `json:"banner"`
	PetaLokalitas string `json:"peta_lokalitas"`
}

// AntarmodaOut (Sub-struct untuk Angkutan Lanjutan)
type AntarmodaOut struct {
	Jenis string   `json:"jenis"`
	Rute  []string `json:"rute"`
}

// FasilitasOut (Sub-struct untuk Retail dan Fasilitas yang Dikelompokkan)
type FasilitasOut struct {
	ID    string `json:"id"`
	Nama  string `json:"nama"`
	Cover string `json:"cover"`
	Tipe  string `json:"tipe"`
}
