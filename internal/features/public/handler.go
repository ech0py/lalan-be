package public

import (
	"log"
	"net/http"

	"lalan-be/internal/response"
	"lalan-be/pkg/message"
)

// Struct untuk handler public.
type PublicHandler struct {
	service PublicService
}

// Fungsi untuk dapatkan daftar kategori.
func (h *PublicHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCategories: received request")
	// Cek method GET
	if r.Method != http.MethodGet {
		response.BadRequest(w, message.MsgNotAllowed)
		return
	}

	categories, err := h.service.GetListCategory()
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, categories, message.MsgSuccess)
}

// Fungsi untuk membuat handler public.
func NewPublicHandler(s PublicService) *PublicHandler {
	return &PublicHandler{service: s}
}
