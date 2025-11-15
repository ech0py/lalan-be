package hoster

import (
	"database/sql"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
	Struktur untuk repositori hoster.

Menyediakan akses ke operasi database untuk hoster.
*/
type hosterRespository struct {
	db *sqlx.DB
}

/*
	Membuat hoster baru di database.

ID dan timestamp hoster dikembalikan setelah penyisipan.
*/
func (r *hosterRespository) CreateHoster(hoster *model.HosterModel) error {
	query := `
		INSERT INTO hoster (
			full_name,
			store_name,
			address,
			phone_number,
			email,
			password_hash,
			profile_photo,
			description,
			tiktok,
			instagram,
			website,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, hoster.FullName, hoster.StoreName, hoster.Address, hoster.PhoneNumber, hoster.Email, hoster.PasswordHash, hoster.ProfilePhoto, hoster.Description, hoster.Tiktok, hoster.Instagram, hoster.Website, hoster.CreatedAt, hoster.UpdatedAt).Scan(&hoster.ID, &hoster.CreatedAt, &hoster.UpdatedAt)
	log.Printf("CreateHoster: inserted hoster with email %s, ID %s", hoster.Email, hoster.ID)
	return err
}

/*
	Mencari hoster berdasarkan email untuk login.

Model hoster dikembalikan jika ditemukan.
*/
func (r *hosterRespository) FindByEmailHosterForLogin(email string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	query := `
		SELECT
			id,
			full_name,
			store_name,
			phone_number,
			email,
			password_hash,
			profile_photo,
			description,
			tiktok,
			instagram,
			website,
			created_at,
			updated_at
		FROM hoster
		WHERE email = $1
	`
	err := r.db.Get(&hoster, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("FindByEmailHosterForLogin: no hoster found for email %s", email)
			return nil, nil
		}
		log.Printf("FindByEmailHosterForLogin: error querying email %s: %v", email, err)
		return nil, err
	}
	log.Printf("FindByEmailHosterForLogin: found hoster for email %s", email)
	return &hoster, nil
}

/*
	Mengambil detail hoster berdasarkan ID.

Model hoster dikembalikan jika ditemukan.
*/
func (r *hosterRespository) GetDetailHoster(id string) (*model.HosterModel, error) {
	var hoster model.HosterModel
	query := `
        SELECT 
            id,
            full_name,
            store_name,
            address,
            phone_number,
            email,
            password_hash,
            profile_photo,
            description,
            tiktok,
            instagram,
            website,
            created_at,
            updated_at
        FROM hoster
        WHERE id = $1
    `
	err := r.db.Get(&hoster, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetDetailHoster: no hoster found for id %s", id)
			return nil, nil
		}
		log.Printf("GetDetailHoster: error for id %s: %v", id, err)
		return nil, err
	}
	log.Printf("GetDetailHoster: found hoster id %s", id)
	return &hoster, nil
}

/*
	Antarmuka untuk operasi repositori hoster.

Mendefinisikan metode untuk CRUD hoster.
*/
type HosterRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmailHosterForLogin(email string) (*model.HosterModel, error)
	GetDetailHoster(id string) (*model.HosterModel, error)
}

/*
	Membuat instance baru dari HosterRepository.

Instance repositori dikembalikan.
*/
func NewHosterRepository(db *sqlx.DB) HosterRepository {
	return &hosterRespository{db: db}
}
