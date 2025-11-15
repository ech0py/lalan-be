package admin

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"

	"lalan-be/internal/model"
)

/*
	Struktur untuk repositori admin.

Menyediakan akses ke operasi database untuk admin.
*/
type adminRepository struct {
	db *sqlx.DB
}

/*
	Membuat admin baru di database.

ID dan timestamp admin dikembalikan setelah penyisipan.
*/
func (r *adminRepository) CreateAdmin(admin *model.AdminModel) error {
	query := `
		INSERT INTO admin (
			email,
			password_hash,
			full_name,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, admin.Email, admin.PasswordHash, admin.FullName, admin.CreatedAt, admin.UpdatedAt).Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt)
	log.Printf("CreateAdmin: inserted admin with email %s, ID %s", admin.Email, admin.ID)
	return err
}

/*
	Mencari admin berdasarkan email untuk login.

Model admin dikembalikan jika ditemukan.
*/
func (r *adminRepository) FindByEmailAdminForLogin(email string) (*model.AdminModel, error) {
	var admin model.AdminModel
	query := `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			created_at,
			updated_at
		FROM admin
		WHERE email = $1
	`
	err := r.db.Get(&admin, query, email)
	// Cek jika tidak ada row
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindByEmailAdminForLogin: no admin found for email %s", email)
			return nil, nil
		}
		log.Printf("FindByEmailAdminForLogin: error querying email %s: %v", email, err)
		return nil, err
	}
	log.Printf("FindByEmailAdminForLogin: found admin for email %s", email)
	return &admin, nil
}

/*
	Membuat kategori baru di database.

ID dan timestamp kategori dikembalikan setelah penyisipan.
*/
func (r *adminRepository) CreateCategory(category *model.CategoryModel) error {
	query := `
		INSERT INTO category (
			name,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	log.Printf("CreateCategory: inserted category with name %s, ID %s", category.Name, category.ID)
	return err
}

/*
	Memperbarui kategori di database.

Kategori berhasil diperbarui atau error dikembalikan.
*/
func (r *adminRepository) UpdateCategory(category *model.CategoryModel) error {
	query := `
		UPDATE category
		SET
			name = $1,
			description = $2,
			updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, category.Name, category.Description, category.UpdatedAt, category.ID)
	log.Printf("UpdateCategory: updated category with ID %s", category.ID)
	return err
}

/*
	Menghapus kategori dari database.

Kategori berhasil dihapus atau error dikembalikan.
*/
func (r *adminRepository) DeleteCategory(id string) error {
	query := `DELETE FROM category WHERE id = $1`
	_, err := r.db.Exec(query, id)
	log.Printf("DeleteCategory: deleted category with ID %s", id)
	return err
}

/*
	Mencari kategori berdasarkan nama.

Model kategori dikembalikan jika ditemukan.
*/
func (r *adminRepository) FindCategoryByName(name string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM category
		WHERE name = $1
	`
	err := r.db.Get(&category, query, name)
	// Cek jika tidak ada row
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindCategoryByName: no category found for name %s", name)
			return nil, nil
		}
		log.Printf("FindCategoryByName: error querying name %s: %v", name, err)
		return nil, err
	}
	log.Printf("FindCategoryByName: found category for name %s", name)
	return &category, nil
}

/*
	Mencari kategori berdasarkan nama kecuali ID.

Model kategori dikembalikan jika ditemukan.
*/
func (r *adminRepository) FindCategoryByNameExceptID(name string, id string) (*model.CategoryModel, error) {
	var category model.CategoryModel
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM category
		WHERE name = $1 AND id != $2
	`
	err := r.db.Get(&category, query, name, id)
	// Cek jika tidak ada row
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindCategoryByNameExceptID: no category found for name %s except ID %s", name, id)
			return nil, nil
		}
		log.Printf("FindCategoryByNameExceptID: error querying name %s except ID %s: %v", name, id, err)
		return nil, err
	}
	log.Printf("FindCategoryByNameExceptID: found category for name %s except ID %s", name, id)
	return &category, nil
}

/*
	Antarmuka untuk operasi repositori admin.

Mendefinisikan metode untuk CRUD admin dan kategori.
*/
type AdminRepository interface {
	CreateAdmin(admin *model.AdminModel) error
	FindByEmailAdminForLogin(email string) (*model.AdminModel, error)
	CreateCategory(category *model.CategoryModel) error
	UpdateCategory(category *model.CategoryModel) error
	DeleteCategory(id string) error
	FindCategoryByName(name string) (*model.CategoryModel, error)
	FindCategoryByNameExceptID(name string, id string) (*model.CategoryModel, error)
}

/*
	Membuat instance baru dari AdminRepository.

Instance repositori dikembalikan.
*/
func NewAdminRepository(db *sqlx.DB) AdminRepository {
	return &adminRepository{db: db}
}
