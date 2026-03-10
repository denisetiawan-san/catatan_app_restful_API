package repository

import (
	"catatan_app/internal/apperror"
	"catatan_app/internal/modul"
	"context"
	"database/sql"
	"errors"
)

// → Struct repository yang memegang koneksi database sebagai dependency
type CatatanRepository struct {
	db *sql.DB
}

// → Constructor untuk inject koneksi database ke repository
func NewCatatanRepository(db *sql.DB) *CatatanRepository {
	return &CatatanRepository{db: db}
}

// → Compile-time check: pastikan CatatanRepository mengimplementasikan CatatanRepo interface
var _ CatatanRepo = (*CatatanRepository)(nil)

// → Menerima context dan objek domain dari service, pintu masuk operasi INSERT ke database
func (r *CatatanRepository) Create(ctx context.Context, note *modul.Catatan) (*modul.Catatan, error) {
	// → [TULIS QUERY SQL] Definisikan query INSERT, pakai ? sebagai placeholder anti SQL injection
	// → arsip default false dan created_at pakai NOW() langsung dari database
	query := `INSERT INTO catatan (judul, isi, arsip, created_at) VALUES (?, ?, false, NOW())`

	// → [EKSEKUSI QUERY] Jalankan query INSERT ke database dengan context
	result, err := r.db.ExecContext(ctx, query, note.Judul, note.Isi)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Query gagal, return error ke service
		return nil, err
	}

	// → [HANDLE HASIL + HANDLE ERROR] Ambil ID baru yang di-generate database setelah INSERT berhasil
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// → [RETURN DATA] Fetch ulang data lengkap dari database berdasarkan ID baru
	// → agar created_at terisi dari database, bukan zero value
	return r.GetByID(ctx, int(id))
}

// → Menerima context, id, dan nilai arsip dari service, pintu masuk operasi UPDATE arsip ke database
func (r *CatatanRepository) SetArsip(ctx context.Context, id int, arsip bool) (*modul.Catatan, error) {
	// → [TULIS QUERY SQL] Definisikan query UPDATE untuk kolom arsip saja
	// → [EKSEKUSI QUERY] Jalankan query UPDATE ke database dengan context
	result, err := r.db.ExecContext(
		ctx,
		"UPDATE catatan SET arsip = ? WHERE id = ?",
		arsip, id,
	)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Query gagal, return error ke service
		return nil, err
	}

	// → [HANDLE HASIL + HANDLE ERROR] Cek berapa baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	// → [HANDLE HASIL + HANDLE ERROR] Jika 0 baris terpengaruh berarti id tidak ditemukan di database
	if rowsAffected == 0 {
		return nil, apperror.ErrNotFound
	}

	// → [RETURN DATA] Fetch ulang data lengkap setelah update agar response selalu fresh dari database
	return r.GetByID(ctx, id)
}

// → Menerima context dan filter arsip dari service, pintu masuk operasi SELECT semua data
func (r *CatatanRepository) GetAll(ctx context.Context, arsip *bool) ([]modul.Catatan, error) {
	// → [TULIS QUERY SQL] Mulai dengan query dasar SELECT semua kolom
	query := `SELECT id, judul, isi, arsip, created_at FROM catatan`
	args := []interface{}{}

	// → [TULIS QUERY SQL] Tambah kondisi WHERE secara dinamis berdasarkan filter arsip
	if arsip != nil {
		// → filter eksplisit dari client, pakai nilai yang dikirim
		query += " WHERE arsip = ?"
		args = append(args, *arsip)
	} else {
		// → tidak ada filter, default tampilkan yang tidak diarsip
		query += " WHERE arsip = false"
	}

	// → [EKSEKUSI QUERY] Jalankan query SELECT ke database dengan context
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Query gagal, return error ke service
		return nil, err
	}
	// → pastikan rows ditutup setelah function selesai agar tidak memory leak
	defer rows.Close()

	// → [HANDLE HASIL + HANDLE ERROR] Iterasi setiap baris hasil query
	var catatan []modul.Catatan
	for rows.Next() {
		var n modul.Catatan
		// → Scan setiap kolom ke field struct modul.Catatan
		if err := rows.Scan(&n.ID, &n.Judul, &n.Isi, &n.Arsip, &n.CreatedAt); err != nil {
			// → [HANDLE HASIL + HANDLE ERROR] Gagal scan baris, return error ke service
			return nil, err
		}
		catatan = append(catatan, n)
	}

	// → [HANDLE HASIL + HANDLE ERROR] Cek error yang terjadi selama iterasi rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// → [RETURN DATA] Pastikan return slice kosong bukan nil jika tidak ada data
	// → agar response JSON tetap [] bukan null
	if catatan == nil {
		catatan = []modul.Catatan{}
	}

	return catatan, nil
}

// → Menerima context dan id dari service, pintu masuk operasi SELECT satu data berdasarkan id
func (r *CatatanRepository) GetByID(ctx context.Context, id int) (*modul.Catatan, error) {
	// → [TULIS QUERY SQL] Definisikan query SELECT dengan kondisi WHERE id
	query := `SELECT id, judul, isi, arsip, created_at FROM catatan WHERE id = ?`

	var n modul.Catatan
	// → [EKSEKUSI QUERY] Jalankan query SELECT dan langsung scan hasilnya ke struct
	// → [HANDLE HASIL + HANDLE ERROR] QueryRowContext hanya return satu baris
	err := r.db.QueryRowContext(ctx, query, id).Scan(&n.ID, &n.Judul, &n.Isi, &n.Arsip, &n.CreatedAt)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Jika tidak ada baris → mapping ke ErrNotFound
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		// → [HANDLE HASIL + HANDLE ERROR] Error lain berarti masalah database
		return nil, err
	}

	// → [RETURN DATA] Return pointer ke struct yang sudah terisi data dari database
	return &n, nil
}

// → Menerima context, id, dan objek domain dari service, pintu masuk operasi UPDATE data
func (r *CatatanRepository) Update(ctx context.Context, id int, catatan *modul.Catatan) (*modul.Catatan, error) {
	// → [TULIS QUERY SQL] Definisikan query UPDATE untuk kolom judul dan isi
	// → [EKSEKUSI QUERY] Jalankan query UPDATE ke database dengan context
	result, err := r.db.ExecContext(
		ctx,
		"UPDATE catatan SET judul = ?, isi = ? WHERE id = ?",
		catatan.Judul, catatan.Isi, id,
	)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Query gagal, return error ke service
		return nil, err
	}

	// → [HANDLE HASIL + HANDLE ERROR] Cek berapa baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	// → [HANDLE HASIL + HANDLE ERROR] Jika 0 baris terpengaruh berarti id tidak ditemukan di database
	if rowsAffected == 0 {
		return nil, apperror.ErrNotFound
	}

	// → [RETURN DATA] Fetch ulang data lengkap setelah update agar response selalu fresh dari database
	return r.GetByID(ctx, id)
}

// → Menerima context dan id dari service, pintu masuk operasi DELETE data dari database
func (r *CatatanRepository) Delete(ctx context.Context, id int) error {
	// → [TULIS QUERY SQL] Definisikan query DELETE berdasarkan id
	// → [EKSEKUSI QUERY] Jalankan query DELETE ke database dengan context
	result, err := r.db.ExecContext(ctx, "DELETE FROM catatan WHERE id = ?", id)
	if err != nil {
		// → [HANDLE HASIL + HANDLE ERROR] Query gagal, return error ke service
		return err
	}

	// → [HANDLE HASIL + HANDLE ERROR] Cek berapa baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// → [HANDLE HASIL + HANDLE ERROR] Jika 0 baris terpengaruh berarti id tidak ditemukan di database
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	// → [RETURN DATA] Hanya return nil, tidak ada data yang dikembalikan karena operasi hapus
	return nil
}

// Setiap function di repository selalu mengikuti flow ini:
// 1. Tulis query SQL    → define query string
// 2. Eksekusi query     → ExecContext atau QueryContext
// 3. Handle hasil       → Scan rows, ambil LastInsertId, cek RowsAffected
// 4. Handle error       → mapping sql.ErrNoRows ke apperror.ErrNotFound
// 5. Return data        → kembalikan model atau error
