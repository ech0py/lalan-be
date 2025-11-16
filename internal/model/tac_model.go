package model

import "time"

/*
Struktur untuk model terms and conditions.
Struktur ini merepresentasikan data terms and conditions dengan field yang diperlukan.
*/
type TermsAndConditionsModel struct {
	ID          string    `json:"id" db:"id"`
	Description []string  `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Foreign key
	UserID string `json:"user_id" db:"user_id"`
}
