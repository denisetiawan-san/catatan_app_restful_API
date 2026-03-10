package service

import (
	"catatan_app/internal/apperror"
	"catatan_app/internal/dto"
	"catatan_app/internal/modul"
	"catatan_app/internal/repository"
	"context"
	"errors"
	"strings"
)

// → Struct service yang memegang repository sebagai dependency via interface
type CatatanService struct {
	repo repository.CatatanRepo
}

// → Constructor untuk inject repository ke service
func NewCatatanService(repo repository.CatatanRepo) *CatatanService {
	return &CatatanService{repo: repo}
}

// → Compile-time check: pastikan CatatanService mengimplementasikan CatatanSvc interface
var _ CatatanSvc = (*CatatanService)(nil)

// → Menerima context dan DTO request dari handler, pintu masuk logic CREATE
func (s *CatatanService) Create(ctx context.Context, req dto.CreateCatatanRequest) (*modul.Catatan, error) {
	// → [TRANSFORMASI DATA] Bersihkan whitespace di awal dan akhir string sebelum diproses
	judul := strings.TrimSpace(req.Judul)
	isi := strings.TrimSpace(req.Isi)

	// → [VALIDASI BISNIS] Aturan bisnis: catatan wajib punya judul, tolak jika kosong setelah trim
	if judul == "" {
		return nil, errors.New("judul harus diisi")
	}

	// → [TRANSFORMASI DATA] Bangun objek domain modul.Catatan dari data yang sudah divalidasi
	catatan := &modul.Catatan{
		Judul: judul,
		Isi:   isi,
	}

	// → [CALL REPOSITORY] Serahkan objek domain ke repository untuk disimpan ke database
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.Create(ctx, catatan)
}

// → Menerima context dan filter arsip dari handler, pintu masuk logic LIST
func (s *CatatanService) List(ctx context.Context, arsip *bool) ([]modul.Catatan, error) {
	// → [CALL REPOSITORY] Tidak ada validasi bisnis khusus, langsung minta data ke repository
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.GetAll(ctx, arsip)
}

// → Menerima context dan id dari handler, pintu masuk logic GET BY ID
func (s *CatatanService) GetByID(ctx context.Context, id int) (*modul.Catatan, error) {
	// → [VALIDASI BISNIS] id tidak boleh kurang dari atau sama dengan 0
	if id <= 0 {
		return nil, apperror.ErrInvalidID
	}

	// → [CALL REPOSITORY] Minta data ke repository berdasarkan id
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.GetByID(ctx, id)
}

// → Menerima context, id, dan DTO request dari handler, pintu masuk logic UPDATE
func (s *CatatanService) Update(ctx context.Context, id int, req dto.UpdateCatatanRequest) (*modul.Catatan, error) {
	// → [VALIDASI BISNIS] id tidak boleh kurang dari atau sama dengan 0
	if id <= 0 {
		return nil, apperror.ErrInvalidID
	}

	// → [CALL REPOSITORY] Fetch data lama dulu sebelum update, sekaligus validasi data exists
	catatan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// → [RETURN HASIL] Jika tidak ditemukan, langsung return error ke handler
		return nil, err
	}

	// → [TRANSFORMASI DATA] Update field hanya jika nilai baru tidak kosong, pertahankan nilai lama jika kosong
	if req.Judul != "" {
		catatan.Judul = strings.TrimSpace(req.Judul)
	}
	if req.Isi != "" {
		catatan.Isi = strings.TrimSpace(req.Isi)
	}

	// → [VALIDASI BISNIS] Pastikan judul tidak kosong setelah proses transformasi
	if catatan.Judul == "" {
		return nil, errors.New("judul harus ada")
	}

	// → [CALL REPOSITORY] Serahkan objek domain yang sudah diupdate ke repository
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.Update(ctx, id, catatan)
}

// → Menerima context dan id dari handler, pintu masuk logic ARSIP
func (s *CatatanService) Arsip(ctx context.Context, id int) (*modul.Catatan, error) {
	// → [VALIDASI BISNIS] id tidak boleh kurang dari atau sama dengan 0
	if id <= 0 {
		return nil, apperror.ErrInvalidID
	}

	// → [CALL REPOSITORY] Set arsip menjadi true di database
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.SetArsip(ctx, id, true)
}

// → Menerima context dan id dari handler, pintu masuk logic UNARSIP
func (s *CatatanService) Unarsip(ctx context.Context, id int) (*modul.Catatan, error) {
	// → [VALIDASI BISNIS] id tidak boleh kurang dari atau sama dengan 0
	if id <= 0 {
		return nil, apperror.ErrInvalidID
	}

	// → [CALL REPOSITORY] Set arsip menjadi false di database
	// → [RETURN HASIL] Hasil dari repository langsung dikembalikan ke handler
	return s.repo.SetArsip(ctx, id, false)
}

// → Menerima context dan id dari handler, pintu masuk logic DELETE
func (s *CatatanService) Delete(ctx context.Context, id int) error {
	// → [VALIDASI BISNIS] id tidak boleh kurang dari atau sama dengan 0
	if id <= 0 {
		return apperror.ErrInvalidID
	}

	// → [CALL REPOSITORY] Serahkan id ke repository untuk dihapus dari database
	// → [RETURN HASIL] Hanya return error, tidak ada data yang dikembalikan
	return s.repo.Delete(ctx, id)
}

// Setiap function di service selalu mengikuti flow ini:
// 1. Validasi bisnis    → cek id > 0, cek field wajib, cek duplikasi
// 2. Transformasi data  → TrimSpace, format ulang, enkripsi password
// 3. Call repository    → minta data atau simpan data ke DB
// 4. Return hasil       → kembalikan model ke handler
