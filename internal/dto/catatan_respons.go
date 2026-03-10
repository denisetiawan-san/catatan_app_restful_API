package dto

import "time"

// → Struct yang mendefinisikan bentuk data yang keluar ke client
// → Hanya field yang boleh dilihat client, field sensitif tidak masuk sini
type CatatanResponse struct {
	ID        int       `json:"id"`         // → tampilkan id ke client
	Judul     string    `json:"judul"`      // → tampilkan judul ke client
	Isi       string    `json:"isi"`        // → tampilkan isi ke client
	Arsip     bool      `json:"arsip"`      // → tampilkan status arsip ke client
	CreatedAt time.Time `json:"created_at"` // → tampilkan waktu dibuat ke client
}

// Pattern-nya selalu:
// 1. Definisi struct nama response
// 2. Field-field yang keluar ke client
// 3. JSON tag untuk mapping
