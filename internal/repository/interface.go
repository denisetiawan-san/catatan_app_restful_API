package repository

import (
	"catatan_app/internal/modul"
	"context"
)

// → Kontrak yang harus dipenuhi oleh semua implementasi repository catatan
// → Service hanya boleh depend ke interface ini, bukan ke struct konkretnya
type CatatanRepo interface {
	// → kontrak operasi INSERT — wajib return data lengkap dan error
	Create(ctx context.Context, note *modul.Catatan) (*modul.Catatan, error)

	// → kontrak operasi SELECT semua — wajib support filter arsip opsional
	GetAll(ctx context.Context, arsip *bool) ([]modul.Catatan, error)

	// → kontrak operasi SELECT satu — wajib return data berdasarkan id
	GetByID(ctx context.Context, id int) (*modul.Catatan, error)

	// → kontrak operasi UPDATE — wajib return data terbaru setelah update
	Update(ctx context.Context, id int, catatan *modul.Catatan) (*modul.Catatan, error)

	// → kontrak operasi DELETE — cukup return error, tidak perlu return data
	Delete(ctx context.Context, id int) error

	// → kontrak operasi UPDATE arsip — wajib return data terbaru setelah update
	SetArsip(ctx context.Context, id int, arsip bool) (*modul.Catatan, error)
}

// Pattern-nya selalu:
// 1. Definisi nama interface
// 2. Daftarkan semua method signature
// 3. Setiap method selalu punya parameter ctx context.Context di posisi pertama
// 4. Setiap method selalu return (data, error) atau hanya error
