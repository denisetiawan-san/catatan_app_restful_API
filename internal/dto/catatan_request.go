package dto

// → Struct yang mendefinisikan bentuk data yang masuk dari client saat CREATE
// → Hanya field yang boleh dikirim client, bukan semua field di database
type CreateCatatanRequest struct {
	Judul string `json:"judul"` // → mapping key "judul" dari JSON body
	Isi   string `json:"isi"`   // → mapping key "isi" dari JSON body
}

// → Struct yang mendefinisikan bentuk data yang masuk dari client saat UPDATE
// → Sama dengan Create untuk kasus ini, tapi dipisah agar bisa berkembang sendiri
type UpdateCatatanRequest struct {
	Judul string `json:"judul"` // → mapping key "judul" dari JSON body
	Isi   string `json:"isi"`   // → mapping key "isi" dari JSON body
}

// Pattern-nya selalu:
// 1. Definisi struct nama request
// 2. Field-field yang masuk dari client
// 3. JSON tag untuk mapping
