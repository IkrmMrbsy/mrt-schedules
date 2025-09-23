package response

// APIResponse adalah struct standar untuk balikan (response) API.
// Semua endpoint API akan memakai format ini supaya konsisten.
// - Code    → angka status HTTP (contoh: 200, 400, 500).
// - Message → pesan singkat tentang hasil request.
// - Data    → isi data utama (bisa apa saja: list, object, atau nil).
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
