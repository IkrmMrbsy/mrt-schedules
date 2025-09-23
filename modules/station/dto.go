package station

// StationIn dipakai untuk menampung data dari API eksternal (API MRT).
// Field "Id" dan "Name" akan otomatis diisi dari JSON yang dikirim API.
// Tag `json:"nid"` artinya Id diisi dari field "nid" pada JSON.
// Tag `json:"title"` artinya Name diisi dari field "title" pada JSON.
type StationIn struct {
	Id   string `json:"nid"`
	Name string `json:"title"`
}

// StationOut dipakai untuk mengirim data balik (response) ke API kita sendiri.
// Bedanya dengan StationIn ada di nama field JSON yang kita tentukan.
// Tag `json:"id"` artinya Id dikirim dengan nama "id".
// Tag `json:"name"` artinya Name dikirim dengan nama "name".
type StationOut struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
