package main

import (
	conectdb "catatan_app/internal/conect_db"
	"catatan_app/internal/handler"
	"catatan_app/internal/repository"
	"catatan_app/internal/router"
	"catatan_app/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// → [LOAD ENVIRONMENT VARIABLE] Baca file .env, jika tidak ada pakai system env
	if err := godotenv.Load(); err != nil {
		log.Println("env tidak ditemukan, pakai system env")
	}

	// → [KONEKSI DATABASE] Buat koneksi ke MySQL menggunakan DSN dari environment variable
	db, err := conectdb.New()
	if err != nil {
		// → Jika koneksi gagal, hentikan program karena tidak bisa jalan tanpa database
		log.Fatalf("failed to connect database: %v", err)
	}
	// → Pastikan koneksi database ditutup dengan benar saat program selesai
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	// → [BUAT REPOSITORY] Inject koneksi database ke repository
	repo := repository.NewCatatanRepository(db)

	// → [BUAT SERVICE] Inject repository ke service via interface CatatanRepo
	service := service.NewCatatanService(repo)

	// → [BUAT HANDLER] Inject service ke handler via interface CatatanSvc
	handler := handler.NewCatatanHandler(service)

	// → [DAFTARKAN ROUTE] Buat mux dan daftarkan semua route catatan
	mux := http.NewServeMux()
	router.Register(mux, handler)

	// → [KONFIGURASI SERVER] Definisikan konfigurasi HTTP server
	server := &http.Server{
		Addr:         ":8080",          // → jalankan di port 8080
		Handler:      mux,              // → pakai mux yang sudah terdaftar routenya
		ReadTimeout:  10 * time.Second, // → maksimal 10 detik untuk baca request
		WriteTimeout: 10 * time.Second, // → maksimal 10 detik untuk tulis response
		IdleTimeout:  60 * time.Second, // → maksimal 60 detik koneksi idle
	}

	// → [JALANKAN SERVER] Jalankan server di goroutine terpisah agar tidak blocking
	go func() {
		log.Println("server berjalan di :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// → Jika server error bukan karena shutdown, hentikan program
			log.Fatalf("server error: %v", err)
		}
	}()

	// → [GRACEFUL SHUTDOWN] Siapkan channel untuk menangkap sinyal stop dari OS
	stop := make(chan os.Signal, 1)
	// → Tangkap sinyal Ctrl+C (Interrupt) dan SIGTERM dari sistem
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// → Blok di sini sampai sinyal diterima
	<-stop
	log.Println("shutting down server....")

	// → Beri waktu maksimal 5 detik untuk request yang sedang berjalan agar selesai dulu
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// → [MATIKAN SERVER] Shutdown server dengan graceful, tunggu request aktif selesai
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server gagal mati: %v", err)
	}

	// → Semua request selesai, program keluar dengan bersih
	log.Println("server keluar dengan selamat..")
}

// Pattern-nya selalu:
// 1. Load environment variable
// 2. Koneksi database
// 3. Buat repository (inject db)
// 4. Buat service (inject repository)
// 5. Buat handler (inject service)
// 6. Buat mux dan daftarkan route
// 7. Konfigurasi dan jalankan HTTP server
// 8. Graceful shutdown
