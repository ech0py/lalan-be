package hoster

import (
	"encoding/json"
	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
	"log"
	"net/http"
	"regexp"
	"strings"
)

/*
	Struktur untuk menangani permintaan hoster.

Menyediakan metode untuk operasi CRUD dan autentikasi.
*/
type HosterHandler struct {
	service HosterService
}

/*
	Struktur data untuk permintaan pembuatan hoster.

Data terstruktur untuk validasi dan pemrosesan.
*/
type HosterRequest struct {
	FullName     string `json:"full_name"`
	StoreName    string `json:"store_name"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	ProfilePhoto string `json:"profile_photo"`
	Description  string `json:"description"`
	Tiktok       string `json:"tiktok"`
	Instagram    string `json:"instagram"`
	Website      string `json:"website"`
}

/*
	Struktur data untuk permintaan login hoster.

Data terstruktur untuk autentikasi.
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
	Membuat hoster baru berdasarkan permintaan.

Hoster berhasil dibuat atau error dikembalikan.
*/
func (h *HosterHandler) CreateHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateHoster: received request")
	// Cek method POST
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req HosterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi full name
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateAdmin: full name required")
		response.BadRequest(w, "Full name is required")
		return
	}

	// Validasi email
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateAdmin: email required")
		response.BadRequest(w, "Email is required")
		return
	}

	// Validasi password
	if strings.TrimSpace(req.Password) == "" {
		log.Printf("CreateAdmin: password required")
		response.BadRequest(w, "Password is required")
		return
	}

	input := &model.HosterModel{
		FullName:     req.FullName,
		StoreName:    req.StoreName,
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		PasswordHash: req.Password,
		Address:      req.Address,
		ProfilePhoto: req.ProfilePhoto,
		Description:  req.Description,
		Tiktok:       req.Tiktok,
		Instagram:    req.Instagram,
		Website:      req.Website,
	}

	err := h.service.CreateHoster(input)
	if err != nil {
		log.Printf("CreateHoster: error creating hoster: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, input, message.MsgSuccess)

}

/*
	Mengautentikasi hoster dengan email dan password.

Token akses dan data pengguna dikembalikan jika berhasil.
*/
func (h *HosterHandler) LoginHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginHoster: received request")
	// Cek method POST
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi email dan password
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginHoster: email or password empty")
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Validasi format email
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginHoster: invalid email format: %s", req.Email)
		response.Error(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	resp, err := h.service.LoginHoster(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginHoster: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("LoginHoster: login successful for email %s", req.Email)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    resp.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600,
	})

	userData := map[string]interface{}{
		"id":            resp.ID,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"token_type":    resp.TokenType,
		"expires_in":    resp.ExpiresIn,
	}

	response.Success(w, 200, userData, "Login successful")
}

/*
	Mengambil detail hoster dari konteks.

Detail hoster dikembalikan jika ditemukan.
*/
func (h *HosterHandler) GetDetailHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailHoster: received request")
	// Cek method GET
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	// Ambil context dari request (asumsikan middleware sudah menyimpan claims)
	ctx := r.Context()

	hoster, err := h.service.GetDetailHoster(ctx)
	if err != nil {
		log.Printf("GetDetailHoster: error getting hoster: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if hoster == nil {
		log.Printf("GetDetailHoster: hoster not found")
		response.Error(w, http.StatusNotFound, "Hoster not found")
		return
	}

	log.Printf("GetDetailHoster: retrieved hoster for ID %s", hoster.ID)
	response.Success(w, http.StatusOK, hoster, "Hoster details retrieved successfully")
}

/*
	Membuat instance baru dari HosterHandler.

Instance HosterHandler dikembalikan.
*/
func NewHosterHandler(s HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}
