package handler

import (
	"catatan_app/internal/apperror"
	"catatan_app/internal/dto"
	"catatan_app/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// → Struct handler yang memegang service sebagai dependency via interface
type CatatanHandler struct {
	service service.CatatanSvc
}

// → Constructor untuk inject service ke handler
func NewCatatanHandler(s service.CatatanSvc) *CatatanHandler {
	return &CatatanHandler{service: s}
}

// → Helper: tulis response JSON ke client dengan status code yang ditentukan
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// → Helper: tulis response error ke client dalam format JSON
func writeError(w http.ResponseWriter, status int, pesan string) {
	writeJSON(w, status, map[string]string{
		"error": pesan,
	})
}

// → Helper: mapping jenis error ke HTTP status yang tepat
func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		// → ErrNotFound → 404 Not Found
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, apperror.ErrInvalidID):
		// → ErrInvalidID → 400 Bad Request
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperror.ErrBadRequest):
		// → ErrBadRequest → 400 Bad Request
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		// → Error tidak dikenal → 500 Internal Server Error
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan pada server")
	}
}

// → Menerima HTTP request dan response writer, pintu masuk request CREATE dari client
func (h *CatatanHandler) Create(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil JSON dari body request, konversi ke struct CreateCatatanRequest
	var req dto.CreateCatatanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// → [VALIDASI INPUT] JSON rusak atau format salah, tolak langsung dengan 400
		writeError(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan ke service, handler tidak proses apapun sendiri
	catatan, err := h.service.Create(r.Context(), req)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi model ke DTO response, kirim 201 Created
	writeJSON(w, http.StatusCreated, dto.ToCatatanResponse(catatan))
}

// → Menerima HTTP request dan response writer, pintu masuk request LIST dari client
func (h *CatatanHandler) List(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil query param "arsip" dari URL jika ada, sifatnya opsional
	var arsip *bool
	if q := r.URL.Query().Get("arsip"); q != "" {
		val, err := strconv.ParseBool(q)
		if err != nil {
			// → [VALIDASI INPUT] Nilai arsip bukan true/false, tolak dengan 400
			writeError(w, http.StatusBadRequest, "nilai arsip tidak valid, gunakan true atau false")
			return
		}
		// → [PARSE INPUT] Simpan nilai arsip sebagai pointer bool
		arsip = &val
	}

	// → [CALL SERVICE] Serahkan ke service beserta filter arsip yang sudah diparsing
	notes, err := h.service.List(r.Context(), arsip)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi slice model ke slice DTO response, kirim 200 OK
	writeJSON(w, http.StatusOK, dto.ToCatatanResponses(notes))
}

// → Menerima HTTP request dan response writer, pintu masuk request GET BY ID dari client
func (h *CatatanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil id dari path URL, konversi string ke integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// → [VALIDASI INPUT] id bukan angka valid, tolak dengan 400
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan id ke service untuk dicari di database
	catatan, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi model ke DTO response, kirim 200 OK
	writeJSON(w, http.StatusOK, dto.ToCatatanResponse(catatan))
}

// → Menerima HTTP request dan response writer, pintu masuk request UPDATE dari client
func (h *CatatanHandler) Update(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil id dari path URL, konversi string ke integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// → [VALIDASI INPUT] id bukan angka valid, tolak dengan 400
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	// → [PARSE INPUT] Ambil JSON dari body request, konversi ke struct UpdateCatatanRequest
	var req dto.UpdateCatatanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// → [VALIDASI INPUT] JSON rusak atau format salah, tolak dengan 400
		writeError(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan id dan request ke service untuk diproses
	catatan, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi model ke DTO response, kirim 200 OK
	writeJSON(w, http.StatusOK, dto.ToCatatanResponse(catatan))
}

// → Menerima HTTP request dan response writer, pintu masuk request ARSIP dari client
func (h *CatatanHandler) Arsip(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil id dari path URL, konversi string ke integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// → [VALIDASI INPUT] id bukan angka valid, tolak dengan 400
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan id ke service untuk diarsipkan
	catatan, err := h.service.Arsip(r.Context(), id)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi model ke DTO response, kirim 200 OK
	writeJSON(w, http.StatusOK, dto.ToCatatanResponse(catatan))
}

// → Menerima HTTP request dan response writer, pintu masuk request UNARSIP dari client
func (h *CatatanHandler) Unarsip(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil id dari path URL, konversi string ke integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// → [VALIDASI INPUT] id bukan angka valid, tolak dengan 400
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan id ke service untuk dikembalikan dari arsip
	catatan, err := h.service.Unarsip(r.Context(), id)
	if err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [BUILD RESPONSE + RETURN JSON] Konversi model ke DTO response, kirim 200 OK
	writeJSON(w, http.StatusOK, dto.ToCatatanResponse(catatan))
}

// → Menerima HTTP request dan response writer, pintu masuk request DELETE dari client
func (h *CatatanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// → [PARSE INPUT] Ambil id dari path URL, konversi string ke integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		// → [VALIDASI INPUT] id bukan angka valid, tolak dengan 400
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	// → [CALL SERVICE] Serahkan id ke service untuk dihapus
	if err := h.service.Delete(r.Context(), id); err != nil {
		// → [HANDLE ERROR] Mapping error dari service ke HTTP status yang tepat
		handleError(w, err)
		return
	}

	// → [RETURN] Tidak ada body response, kirim 204 No Content sebagai tanda berhasil dihapus
	w.WriteHeader(http.StatusNoContent)
}

// Setiap function di handler selalu mengikuti flow ini:
// 1. Parse input        → ambil id dari URL, decode JSON body
// 2. Validasi input     → cek format, bukan logik bisnis
// 3. Call service       → serahkan ke service, jangan proses sendiri
// 4. Handle error       → mapping error ke HTTP status yang tepat
// 5. Build response     → konversi model ke DTO response
// 6. Return JSON        → tulis ke http.ResponseWriter
