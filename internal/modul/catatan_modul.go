package modul

import "time"

// → Definisi struct domain yang merepresentasikan tabel catatan di database
// → Tidak ada JSON tag karena ini murni domain object, bukan response API
type Catatan struct {
	ID        int       // → representasi kolom id di tabel
	Judul     string    // → representasi kolom judul di tabel
	Isi       string    // → representasi kolom isi di tabel
	Arsip     bool      // → representasi kolom arsip di tabel
	CreatedAt time.Time // → representasi kolom created_at di tabel
}

// Pattern-nya selalu:
// 1. Definisi struct nama domain
// 2. Field-field yang merepresentasikan kolom tabel database
// 3. Tidak ada JSON tag
// 4. Tidak ada logic apapun
