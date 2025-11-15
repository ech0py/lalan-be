package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"lalan-be/internal/model"
	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

// Struct untuk handler admin.
type AdminHandler struct {
	service AdminService
}

// Struct untuk request admin.
type AdminRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct untuk request login admin.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct untuk request kategori.
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Fungsi untuk membuat admin baru.
func (h *AdminHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateAdmin: received request")
	// Cek method POST
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}
	var req AdminRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateAdmin: invalid JSON: %v", err)
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

	input := &model.AdminModel{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	err := h.service.CreateAdmin(input)
	if err != nil {
		log.Printf("CreateAdmin: error creating admin: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, input, message.MsgSuccess)
}

// Fungsi untuk login admin.
func (h *AdminHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginAdmin: received request")
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
		log.Printf("LoginAdmin: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi email dan password
	if req.Email == "" || req.Password == "" {
		log.Printf("LoginAdmin: email or password empty")
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Validasi format email
	if !emailRegex.MatchString(req.Email) {
		log.Printf("LoginAdmin: invalid email format: %s", req.Email)
		response.Error(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	resp, err := h.service.LoginAdmin(req.Email, req.Password)
	if err != nil {
		log.Printf("LoginAdmin: login failed: %v", err)
		response.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	log.Printf("LoginAdmin: login successful for email %s", req.Email)
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

// Fungsi untuk membuat kategori baru.
func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateCategory: received request")
	// Cek method POST
	if r.Method != http.MethodPost {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("CreateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("CreateCategory: name required")
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("CreateCategory: name too long")
		response.BadRequest(w, message.MsgCategoryNameTooLong)
		return
	}

	input := &model.CategoryModel{
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.CreateCategory(input)
	if err != nil {
		log.Printf("CreateCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, input, message.MsgCategoryCreatedSuccess)
}

// Fungsi untuk update kategori.
func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("UpdateCategory: received request")
	// Cek method PUT
	if r.Method != http.MethodPut {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("UpdateCategory: id required")
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}

	var req CategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	// Decode JSON
	if err := decoder.Decode(&req); err != nil {
		log.Printf("UpdateCategory: invalid JSON: %v", err)
		response.BadRequest(w, message.MsgBadRequest)
		return
	}

	// Validasi name
	if strings.TrimSpace(req.Name) == "" {
		log.Printf("UpdateCategory: name required")
		response.BadRequest(w, message.MsgCategoryNameRequired)
		return
	}

	// Validasi panjang name
	if len(req.Name) > 255 {
		log.Printf("UpdateCategory: name too long")
		response.BadRequest(w, message.MsgCategoryNameTooLong)
		return
	}

	input := &model.CategoryModel{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.service.UpdateCategory(input)
	if err != nil {
		log.Printf("UpdateCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, input, message.MsgCategoryUpdatedSuccess)
}

// Fungsi untuk hapus kategori.
func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteCategory: received request")
	// Cek method DELETE
	if r.Method != http.MethodDelete {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	// Validasi ID
	if strings.TrimSpace(id) == "" {
		log.Printf("DeleteCategory: id required")
		response.BadRequest(w, message.MsgCategoryIDRequired)
		return
	}

	err := h.service.DeleteCategory(id)
	if err != nil {
		log.Printf("DeleteCategory: error: %v", err)
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, nil, message.MsgCategoryDeletedSuccess)
}

// Fungsi untuk membuat handler admin.
func NewAdminHandler(s AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}
