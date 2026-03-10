package apperror

import "errors"

// → Definisi sentinel error terpusat yang dipakai oleh semua layer
// → Tujuannya agar handler bisa mapping jenis error ke HTTP status yang tepat
var (
	// → Dipakai ketika data tidak ditemukan di database → handler mapping ke 404
	ErrNotFound = errors.New("catatan tidak ditemukan")

	// → Dipakai ketika id yang dikirim tidak valid (misal <= 0) → handler mapping ke 400
	ErrInvalidID = errors.New("id tidak valid")

	// → Dipakai ketika request tidak memenuhi syarat → handler mapping ke 400
	ErrBadRequest = errors.New("request tidak valid")
)
