package public

import (
	"log"
	"net/http"

	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

/*
Struktur untuk handler public.
Struktur ini menangani permintaan publik.
*/
type PublicHandler struct {
	service PublicService
}

/*
Metode untuk mendapatkan semua kategori.
Daftar kategori dikembalikan.
*/
func (h *PublicHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCategories: received request")
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	categories, err := h.service.GetAllCategory()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

/*
Metode untuk mendapatkan semua item.
Daftar item dikembalikan.
*/
func (h *PublicHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
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
Metode untuk mendapatkan semua syarat dan ketentuan.
Daftar syarat dan ketentuan dikembalikan.
*/
func (h *PublicHandler) GetAllTermsAndConditions(w http.ResponseWriter, r *http.Request) {
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
Fungsi untuk membuat instance baru dari PublicHandler.
Instance handler dikembalikan.
*/
func NewPublicHandler(s PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
