package dto

import "catatan_app/internal/modul"

// → [KONVERSI SINGLE] Konversi satu objek domain *modul.Catatan ke CatatanResponse
// → Dipakai oleh handler untuk operasi yang return satu data: Create, GetByID, Update, Arsip, Unarsip
func ToCatatanResponse(catatan *modul.Catatan) CatatanResponse {
	// → Mapping setiap field dari domain object ke response struct
	return CatatanResponse{
		ID:        catatan.ID,
		Judul:     catatan.Judul,
		Isi:       catatan.Isi,
		Arsip:     catatan.Arsip,
		CreatedAt: catatan.CreatedAt,
	}
}

// → [KONVERSI SLICE] Konversi slice []modul.Catatan ke slice []CatatanResponse
// → Dipakai oleh handler untuk operasi yang return banyak data: List
func ToCatatanResponses(notes []modul.Catatan) []CatatanResponse {
	// → Buat slice response dengan kapasitas sama dengan slice input
	responses := make([]CatatanResponse, len(notes))
	for i, catatan := range notes {
		// → Konversi setiap item menggunakan fungsi single di atas
		responses[i] = ToCatatanResponse(&catatan)
	}
	return responses
}

// Pattern-nya selalu:
// 1. Function konversi single → dari *model ke Response
// 2. Function konversi slice  → dari []model ke []Response
