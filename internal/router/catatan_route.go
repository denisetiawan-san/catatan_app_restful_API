package router

import (
	"catatan_app/internal/handler"
	"net/http"
)

// → Helper: tulis response 405 Method Not Allowed dalam format JSON
func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error":"method not allowed"}`))
}

// → [DEFINISI FUNCTION REGISTER] Daftarkan semua route resource catatan ke mux
func Register(mux *http.ServeMux, h *handler.CatatanHandler) {

	// → [DAFTARKAN ENDPOINT] Path /catatan untuk operasi koleksi
	mux.HandleFunc("/catatan", func(w http.ResponseWriter, r *http.Request) {
		// → [SWITCH METHOD HTTP] Arahkan ke handler yang sesuai berdasarkan method
		switch r.Method {
		case http.MethodPost:
			// → POST /catatan → handler Create
			h.Create(w, r)
		case http.MethodGet:
			// → GET /catatan → handler List
			h.List(w, r)
		default:
			// → Method selain POST dan GET → tolak dengan 405
			methodNotAllowed(w)
		}
	})

	// → [DAFTARKAN ENDPOINT] Path /catatan/{id} untuk operasi per item
	mux.HandleFunc("/catatan/{id}", func(w http.ResponseWriter, r *http.Request) {
		// → [SWITCH METHOD HTTP] Arahkan ke handler yang sesuai berdasarkan method
		switch r.Method {
		case http.MethodGet:
			// → GET /catatan/{id} → handler GetByID
			h.GetByID(w, r)
		case http.MethodPut:
			// → PUT /catatan/{id} → handler Update
			h.Update(w, r)
		case http.MethodDelete:
			// → DELETE /catatan/{id} → handler Delete
			h.Delete(w, r)
		default:
			// → Method selain GET, PUT, DELETE → tolak dengan 405
			methodNotAllowed(w)
		}
	})

	// → [DAFTARKAN ENDPOINT] Path /catatan/{id}/arsip untuk operasi arsip
	mux.HandleFunc("/catatan/{id}/arsip", func(w http.ResponseWriter, r *http.Request) {
		// → [SWITCH METHOD HTTP] Hanya terima PATCH, selain itu tolak
		if r.Method != http.MethodPatch {
			methodNotAllowed(w)
			return
		}
		// → PATCH /catatan/{id}/arsip → handler Arsip
		h.Arsip(w, r)
	})

	// → [DAFTARKAN ENDPOINT] Path /catatan/{id}/unarsip untuk operasi unarsip
	mux.HandleFunc("/catatan/{id}/unarsip", func(w http.ResponseWriter, r *http.Request) {
		// → [SWITCH METHOD HTTP] Hanya terima PATCH, selain itu tolak
		if r.Method != http.MethodPatch {
			methodNotAllowed(w)
			return
		}
		// → PATCH /catatan/{id}/unarsip → handler Unarsip
		h.Unarsip(w, r)
	})
}

// Pattern-nya selalu:
// 1. Definisi function Register untuk resource tersebut
// 2. Daftarkan endpoint dengan path-nya ->mux.Handlefunc
// 3. Switch method HTTP → arahkan ke handler yang sesuai
// 4. Default case selalu methodNotAllowed
