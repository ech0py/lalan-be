package response

import (
	"encoding/json"
	"net/http"
)

/*
Struktur untuk respons API.
Struktur ini digunakan untuk format respons JSON standar.
*/
type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

/*
Fungsi untuk mengirim respons sukses.
Respons JSON dengan status sukses dikirim.
*/
func Success(w http.ResponseWriter, code int, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Data:    data,
		Message: message,
		Success: true,
	})
}

/*
Fungsi untuk mengirim respons error.
Respons JSON dengan status error dikirim.
*/
func Error(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
		Success: false,
	})
}

/*
Fungsi untuk mengirim respons OK.
Respons JSON dengan status OK dikirim.
*/
func OK(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusOK, data, msg)
}

/*
Fungsi untuk mengirim respons Created.
Respons JSON dengan status Created dikirim.
*/
func Created(w http.ResponseWriter, data any, msg string) {
	Success(w, http.StatusCreated, data, msg)
}

/*
Fungsi untuk mengirim respons BadRequest.
Respons JSON dengan status BadRequest dikirim.
*/
func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}

/*
Fungsi untuk mengirim respons Unauthorized.
Respons JSON dengan status Unauthorized dikirim.
*/
func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, msg)
}

/*
Fungsi untuk mengirim respons Forbidden.
Respons JSON dengan status Forbidden dikirim.
*/
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}
