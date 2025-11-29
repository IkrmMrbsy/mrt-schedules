package station

// StationIn dipakai untuk menampung data dari API eksternal (API MRT).
// Field "Id" dan "Name" akan otomatis diisi dari JSON yang dikirim API.
// Tag `json:"nid"` artinya Id diisi dari field "nid" pada JSON.
// Tag `json:"title"` artinya Name diisi dari field "title" pada JSON.
type StationIn struct {
	ID          string `json:"nid"`
	NamaStasiun string `json:"title"`

	Antarmoda     string        `json:"antarmodas"`
	PetaLokalitas string        `json:"peta_lokalitas"`
	Banner        string        `json:"banner"`
	Retails       []RetailIn    `json:"retails"`
	Fasilitas     []FasilitasIn `json:"fasilitas"`
}

type RetailIn struct {
	ID          string `json:"nid"`
	Judul       string `json:"title"`
	JenisRetail string `json:"jenis_retail"`
	Cover       string `json:"cover"`
}

type FasilitasIn struct {
	ID             string `json:"nid"`
	Judul          string `json:"title"`
	JenisFasilitas string `json:"jenis_fasilitas"`
	Cover          string `json:"cover"`
}

type ScheduleIn struct {
	IDStasiun             string `json:"nid"`
	NamaStasiun           string `json:"title"`
	JadwalBundaranHIBiasa string `json:"jadwal_hi_biasa"`
	JadwalBundaranHILibur string `json:"jadwal_hi_libur"`
	JadwalLebakBulusBiasa string `json:"jadwal_lb_biasa"`
	JadwalLebakBulusLibur string `json:"jadwal_lb_libur"`
}

type EstimasiIn struct {
	IDStasiunTujuan string `json:"stasiun_nid"`
	Tarif           string `json:"tarif"`
	Waktu           string `json:"waktu"`
}

type FareIn struct {
	ID       string       `json:"nid"`
	Nama     string       `json:"title"`
	Tarif    string       `json:"tarif"`
	Estimasi []EstimasiIn `json:"estimasi"`
}
