package hoster

import (
	"database/sql"
	"encoding/json"
	"lalan-be/internal/model"
	"log"

	"github.com/jmoiron/sqlx"
)

/*
Struktur untuk repositori hoster.
Struktur ini menyediakan akses ke operasi database untuk hoster.
*/
type hosterRespository struct {
	db *sqlx.DB
}

/*
Metode untuk membuat hoster baru di database.
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
Metode untuk mencari hoster berdasarkan email untuk login.
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
Metode untuk mengambil detail hoster berdasarkan ID.
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
Metode untuk membuat item baru di database.
Item berhasil dibuat atau error dikembalikan.
*/
func (r *hosterRespository) CreateItem(item *model.ItemModel) error {
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("CreateItem: error marshaling photos: %v", err)
		return err
	}
	query := `
		INSERT INTO item (
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`
	_, err = r.db.Exec(query, item.ID, item.Name, item.Description, photosJSON,
		item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount,
		item.CategoryID, item.UserID)
	if err != nil {
		log.Printf("CreateItem: error inserting item: %v", err)
		return err
	}
	return nil
}

/*
Metode untuk mencari item berdasarkan ID.
Model item dikembalikan jika ditemukan.
*/
func (r *hosterRespository) FindItemNameByID(id string) (*model.ItemModel, error) {
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
		WHERE id = $1
		LIMIT 1
	`
	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindByID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
Metode untuk mencari item berdasarkan nama dan user ID.
Model item dikembalikan jika ditemukan.
*/
func (r *hosterRespository) FindItemNameByUserID(name string, userId string) (*model.ItemModel, error) {
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
		WHERE name = $1 AND user_id = $2
		LIMIT 1
	`
	var item model.ItemModel
	var photosJSON []byte
	err := r.db.QueryRow(query, name, userId).Scan(
		&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock,
		&item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount,
		&item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindItemNameByUserID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
		log.Printf("Unmarshal photos error: %v", err)
		return nil, err
	}

	return &item, nil
}

/*
Metode untuk mengambil semua item dari database.
Daftar model item dikembalikan.
*/
func (r *hosterRespository) GetAllItems() ([]*model.ItemModel, error) {
	query := `
		SELECT
			id,
			name,
			description,
			photos,
			stock,
			pickup_type,
			price_per_day,
			deposit,
			discount,
			category_id,
			user_id,
			created_at,
			updated_at
		FROM item
	`
	var items []*model.ItemModel
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ItemModel
		var photosJSON []byte
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &photosJSON, &item.Stock, &item.PickupType, &item.PricePerDay, &item.Deposit, &item.Discount, &item.CategoryID, &item.UserID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(photosJSON, &item.Photos); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

/*
Metode untuk memperbarui item di database.
Item diperbarui berdasarkan ID.
*/
func (r *hosterRespository) UpdateItem(item *model.ItemModel) error {
	query := `
		UPDATE item
		SET
			name = $1,
			description = $2,
			photos = $3,
			stock = $4,
			pickup_type = $5,
			price_per_day = $6,
			deposit = $7,
			discount = $8,
			category_id = $9,
			updated_at = $10
		WHERE id = $11
	`
	photosJSON, err := json.Marshal(item.Photos)
	if err != nil {
		log.Printf("UpdateItem: error marshaling photos: %v", err)
		return err
	}
	_, err = r.db.Exec(query, item.Name, item.Description, photosJSON, item.Stock, item.PickupType, item.PricePerDay, item.Deposit, item.Discount, item.CategoryID, item.UpdatedAt, item.ID)
	if err != nil {
		log.Printf("UpdateItem: error updating item: %v", err)
		return err
	}
	return nil
}

/*
Metode untuk menghapus item dari database.
Item dihapus berdasarkan ID.
*/
func (r *hosterRespository) DeleteItem(id string) error {
	query := `DELETE FROM item WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteItem: error deleting item: %v", err)
		return err
	}
	return nil
}

/*
Metode untuk membuat terms and conditions baru di database.
Terms and conditions berhasil dibuat atau error dikembalikan.
*/
func (r *hosterRespository) CreateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error marshaling description: %v", err)
		return err
	}
	query := `
		INSERT INTO tnc (
			id,
			user_id,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, NOW(), NOW())
	`
	_, err = r.db.Exec(query, tac.ID, tac.UserID, descriptionJSON)
	if err != nil {
		log.Printf("CreateTermsAndConditions: error inserting tnc: %v", err)
		return err
	}
	return nil
}

/*
Metode untuk mencari terms and conditions berdasarkan ID.
Model terms and conditions dikembalikan jika ditemukan.
*/
func (r *hosterRespository) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	query := `
		SELECT
			id,
			user_id,
			description,
			created_at,
			updated_at
		FROM tnc
		WHERE id = $1
		LIMIT 1
	`
	var tac model.TermsAndConditionsModel
	var descriptionJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&tac.ID, &tac.UserID, &descriptionJSON, &tac.CreatedAt, &tac.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("FindTermByID error: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(descriptionJSON, &tac.Description); err != nil {
		log.Printf("Unmarshal description error: %v", err)
		return nil, err
	}

	return &tac, nil
}

/*
Metode untuk mengambil semua terms and conditions dari database.
Daftar model terms and conditions dikembalikan.
*/
func (r *hosterRespository) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	query := `
		SELECT
			id,
			user_id,
			description,
			created_at,
			updated_at
		FROM tnc
	`
	var terms []*model.TermsAndConditionsModel
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tac model.TermsAndConditionsModel
		var descriptionJSON []byte
		err := rows.Scan(&tac.ID, &tac.UserID, &descriptionJSON, &tac.CreatedAt, &tac.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(descriptionJSON, &tac.Description); err != nil {
			return nil, err
		}
		terms = append(terms, &tac)
	}
	return terms, nil
}

/*
Metode untuk memperbarui terms and conditions di database.
Terms and conditions diperbarui berdasarkan ID.
*/
func (r *hosterRespository) UpdateTermsAndConditions(tac *model.TermsAndConditionsModel) error {
	descriptionJSON, err := json.Marshal(tac.Description)
	if err != nil {
		log.Printf("UpdateTerm: error marshaling description: %v", err)
		return err
	}
	query := `
		UPDATE tnc
		SET
			description = $1,
			updated_at = NOW()
		WHERE id = $2
	`
	_, err = r.db.Exec(query, descriptionJSON, tac.ID)
	if err != nil {
		log.Printf("UpdateTerm: error updating tnc: %v", err)
		return err
	}
	return nil
}

/*
Metode untuk menghapus terms and conditions dari database.
Terms and conditions dihapus berdasarkan ID.
*/
func (r *hosterRespository) DeleteTermsAndConditions(id string) error {
	query := `DELETE FROM tnc WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("DeleteTerm: error deleting tnc: %v", err)
		return err
	}
	return nil
}

/*
Antarmuka untuk operasi repositori hoster.
Mendefinisikan metode untuk CRUD hoster.
*/
type HosterRepository interface {
	CreateHoster(hoster *model.HosterModel) error
	FindByEmailHosterForLogin(email string) (*model.HosterModel, error)
	GetDetailHoster(id string) (*model.HosterModel, error)
	CreateItem(item *model.ItemModel) error
	FindItemNameByUserID(name string, userId string) (*model.ItemModel, error)
	FindItemNameByID(id string) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	UpdateItem(item *model.ItemModel) error
	DeleteItem(id string) error
	CreateTermsAndConditions(tac *model.TermsAndConditionsModel) error
	FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
	UpdateTermsAndConditions(tac *model.TermsAndConditionsModel) error
	DeleteTermsAndConditions(id string) error
}

/*
Fungsi untuk membuat instance baru dari HosterRepository.
Instance repositori dikembalikan.
*/
func NewHosterRepository(db *sqlx.DB) HosterRepository {
	return &hosterRespository{db: db}
}
