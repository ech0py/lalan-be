package public

import (
	"lalan-be/internal/model"
)

// Interface untuk service public.
type PublicService interface {
	GetListCategory() ([]*model.CategoryModel, error)
}

// Struct untuk service public.
type publicService struct {
	repo PublicRepository
}

// Fungsi untuk dapatkan daftar kategori.
func (s *publicService) GetListCategory() ([]*model.CategoryModel, error) {
	return s.repo.GetListCategory()
}

// Fungsi untuk membuat service public.
func NewPublicService(repo PublicRepository) PublicService {
	return &publicService{repo: repo}
}
