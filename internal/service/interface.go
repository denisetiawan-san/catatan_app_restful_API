package service

import (
	"catatan_app/internal/dto"
	"catatan_app/internal/modul"
	"context"
)

// → Kontrak yang harus dipenuhi oleh semua implementasi service catatan
// → Handler hanya boleh depend ke interface ini, bukan ke struct konkretnya
type CatatanSvc interface {
	// → kontrak operasi create — terima DTO request, return domain object dan error
	Create(ctx context.Context, req dto.CreateCatatanRequest) (*modul.Catatan, error)

	// → kontrak operasi list — terima filter arsip opsional, return slice domain object
	List(ctx context.Context, arsip *bool) ([]modul.Catatan, error)

	// → kontrak operasi get by id — terima id, return satu domain object
	GetByID(ctx context.Context, id int) (*modul.Catatan, error)

	// → kontrak operasi update — terima id dan DTO request, return domain object terupdate
	Update(ctx context.Context, id int, req dto.UpdateCatatanRequest) (*modul.Catatan, error)

	// → kontrak operasi arsip — terima id, return domain object dengan arsip true
	Arsip(ctx context.Context, id int) (*modul.Catatan, error)

	// → kontrak operasi unarsip — terima id, return domain object dengan arsip false
	Unarsip(ctx context.Context, id int) (*modul.Catatan, error)

	// → kontrak operasi delete — terima id, cukup return error saja
	Delete(ctx context.Context, id int) error
}

// Pattern-nya sama persis dengan interface repository:
// 1. Definisi nama interface
// 2. Daftarkan semua method signature
// 3. Setiap method selalu punya parameter ctx context.Context di posisi pertama
// 4. Setiap method selalu return (data, error) atau hanya error
