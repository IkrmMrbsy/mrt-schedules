package client

import (
	"errors"
	"io"
	"net/http"
)

// DoRequest adalah fungsi helper untuk melakukan HTTP GET request.
// - Param client: http.Client yang dipakai (sudah ada timeout dll).
// - Param url: alamat tujuan request.
// - Return []byte: isi response dalam bentuk byte.
// - Return error: error kalau ada masalah.
func DoRequest(client *http.Client, url string) ([]byte, error) {

	// Kirim request GET ke URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// Pastikan body response ditutup setelah selesai dipakai
	defer resp.Body.Close()

	// Kalau status code bukan 200 (OK), anggap error
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	// Baca semua isi response body jadi []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Kembalikan isi response
	return body, nil
}
