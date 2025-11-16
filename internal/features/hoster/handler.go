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

	"github.com/gorilla/mux"
)

/*
Struktur untuk menangani permintaan terkait hoster.
Struktur ini menyediakan layanan untuk operasi hoster.
*/
type HosterHandler struct {
	service HosterService
}

/*
Struktur untuk permintaan pembuatan hoster.
Struktur ini berisi data yang diperlukan untuk membuat hoster baru.
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
Struktur untuk permintaan login hoster.
Struktur ini berisi kredensial untuk autentikasi hoster.
*/
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/*
Metode untuk membuat hoster baru.
Metode ini memvalidasi input dan membuat hoster melalui layanan.
*/
func (h *HosterHandler) CreateHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateHoster: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req HosterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if strings.TrimSpace(req.FullName) == "" {
		log.Printf("CreateAdmin: full name required")
		response.BadRequest(w, "Full name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		log.Printf("CreateAdmin: email required")
		response.BadRequest(w, "Email is required")
		return
	}
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
Metode untuk login hoster.
Metode ini memvalidasi kredensial dan mengembalikan token autentikasi.
*/
func (h *HosterHandler) LoginHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginHoster: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("LoginHoster: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginHoster: email or password empty")
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
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
Metode untuk mendapatkan detail hoster.
Metode ini mengambil data hoster berdasarkan konteks permintaan.
*/
func (h *HosterHandler) GetDetailHoster(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetDetailHoster: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
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
Metode untuk membuat item baru.
Metode ini memvalidasi dan membuat item melalui layanan.
*/
func (h *HosterHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateItem: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.CreateItem(ctx, &req)
	if err != nil {
		log.Printf("CreateItem: error creating item: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, item, message.MsgItemCreatedSuccess)
}

/*
Metode untuk mendapatkan item berdasarkan ID.
Metode ini mengambil data item spesifik dari layanan.
*/
func (h *HosterHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetItemByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}
	item, err := h.service.GetItemByID(id)
	if err != nil {
		log.Printf("GetItemByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, item, "Item retrieved successfully")
}

/*
Metode untuk mendapatkan semua item.
Metode ini mengambil daftar semua item dari layanan.
*/
func (h *HosterHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllItems: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	items, err := h.service.GetAllItems()
	if err != nil {
		log.Printf("GetAllItems: error: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, items, "Items retrieved successfully")
}

/*
Metode untuk memperbarui item.
Metode ini memvalidasi dan memperbarui item melalui layanan.
*/
func (h *HosterHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateItem: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}
	var req model.ItemModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateItem: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	item, err := h.service.UpdateItem(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateItem: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, item, message.MsgItemUpdatedSuccess)
}

/*
Metode untuk menghapus item.
Metode ini menghapus item berdasarkan ID melalui layanan.
*/
func (h *HosterHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteItem: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, message.MsgItemIDRequired)
		return
	}
	ctx := r.Context()
	err := h.service.DeleteItem(ctx, id)
	if err != nil {
		log.Printf("DeleteItem: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, nil, message.MsgItemDeletedSuccess)
}

/*
Metode untuk membuat syarat dan ketentuan.
Metode ini memvalidasi dan membuat syarat dan ketentuan melalui layanan.
*/
func (h *HosterHandler) CreateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateTermsAndConditions: received request")
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.CreateTermsAndConditions(ctx, &req)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, tac, "Terms and conditions created successfully")
}

/*
Metode untuk menemukan syarat dan ketentuan berdasarkan ID.
Metode ini mengambil data syarat dan ketentuan spesifik dari layanan.
*/
func (h *HosterHandler) FindTermsAndConditionsByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("FindTermsAndConditionsByID: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if id == "" {
		response.BadRequest(w, "ID is required")
		return
	}
	tac, err := h.service.FindTermsAndConditionsByID(id)
	if err != nil {
		log.Printf("FindTermsAndConditionsByID: error: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, tac, "Terms and conditions retrieved successfully")
}

/*
Metode untuk mendapatkan semua syarat dan ketentuan.
Metode ini mengambil daftar semua syarat dan ketentuan dari layanan.
*/
func (h *HosterHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllTermsAndConditions: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	tacs, err := h.service.GetAllTermsAndConditions()
	if err != nil {
		log.Printf("GetAllTermsAndConditions: error: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, tacs, "Terms and conditions retrieved successfully")
}

/*
Metode untuk memperbarui syarat dan ketentuan.
Metode ini memvalidasi dan memperbarui syarat dan ketentuan melalui layanan.
*/
func (h *HosterHandler) UpdateTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateTermsAndConditions: received request")
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "ID is required")
		return
	}
	var req model.TermsAndConditionsModel
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateTermsAndConditions: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}
	ctx := r.Context()
	tac, err := h.service.UpdateTermsAndConditions(ctx, id, &req)
	if err != nil {
		log.Printf("UpdateTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, tac, "Terms and conditions updated successfully")
}

/*
Metode untuk menghapus syarat dan ketentuan.
Metode ini menghapus syarat dan ketentuan berdasarkan ID melalui layanan.
*/
func (h *HosterHandler) DeleteTermsAndConditions(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteTermsAndConditions: received request")
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if strings.TrimSpace(id) == "" {
		response.BadRequest(w, "ID is required")
		return
	}
	ctx := r.Context()
	err := h.service.DeleteTermsAndConditions(ctx, id)
	if err != nil {
		log.Printf("DeleteTermsAndConditions: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}
	response.OK(w, nil, "Terms and conditions deleted successfully")
}

/*
Fungsi untuk membuat instance baru dari HosterHandler.
Fungsi ini menginisialisasi handler dengan layanan yang diberikan.
*/
func NewHosterHandler(s HosterService) *HosterHandler {
	return &HosterHandler{service: s}
}
