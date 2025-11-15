package public

import (
	"log"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

// Interface untuk repository public.
type PublicRepository interface {
	GetListCategory() ([]*model.CategoryModel, error)
}

// Struct untuk repository public.
type publicRepository struct {
	db *sqlx.DB
}

// Fungsi untuk dapatkan daftar kategori.
func (r *publicRepository) GetListCategory() ([]*model.CategoryModel, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM category ORDER BY created_at DESC"
	var categories []*model.CategoryModel
	err := r.db.Select(&categories, query)
	// Cek error query
	if err != nil {
		log.Printf("GetListCategory error: %v", err)
		return nil, err
	}
	return categories, nil
}

// Fungsi untuk membuat repository public.
func NewPublicRepository(db *sqlx.DB) PublicRepository {
	return &publicRepository{db: db}
}
