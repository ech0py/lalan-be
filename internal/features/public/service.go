package public

import (
	"lalan-be/internal/model"
)

/*
Struktur untuk layanan public.
Struktur ini menyediakan logika bisnis untuk operasi data publik.
*/
type publicService struct {
	repo PublicRepository
}

/*
Metode untuk mendapatkan semua kategori.
Daftar model kategori dikembalikan.
*/
func (s *publicService) GetAllCategory() ([]*model.CategoryModel, error) {
	return s.repo.GetAllCategory()
}

/*
Metode untuk mendapatkan semua item.
Daftar model item dikembalikan.
*/
func (s *publicService) GetAllItems() ([]*model.ItemModel, error) {
	return s.repo.GetAllItems()
}

/*
Metode untuk mendapatkan semua syarat dan ketentuan.
Daftar model syarat dan ketentuan dikembalikan.
*/
func (s *publicService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	return s.repo.GetAllTermsAndConditions()
}

/*
Antarmuka untuk layanan public.
Antarmuka ini mendefinisikan metode untuk operasi data publik.
*/
type PublicService interface {
	GetAllCategory() ([]*model.CategoryModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
}

/*
Fungsi untuk membuat instance baru dari PublicService.
Instance layanan dikembalikan.
*/
func NewPublicService(repo PublicRepository) PublicService {
	return &publicService{repo: repo}
}
