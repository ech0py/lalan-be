package hoster

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lalan-be/internal/config"
	"lalan-be/internal/middleware"
	"lalan-be/internal/model"
	"lalan-be/pkg/message"
)

/*
	Struktur untuk layanan hoster.

Menyediakan logika bisnis untuk operasi hoster.
*/
type hosterService struct {
	repo HosterRepository
}

/*
	Menghasilkan token JWT untuk hoster.

Respons token dikembalikan jika berhasil.
*/
func (s *hosterService) generateTokenHoster(userID string) (*HosterResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "hoster",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	return &HosterResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

/*
	Mengautentikasi hoster dengan email dan password.

Respons token dikembalikan jika berhasil.
*/
func (s *hosterService) LoginHoster(email, password string) (*HosterResponse, error) {
	hoster, err := s.repo.FindByEmailHosterForLogin(email)
	// Cek error atau hoster tidak ada
	if err != nil || hoster == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verifikasi password
	if bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenHoster(hoster.ID)
}

/*
	Membuat hoster baru dengan hashing password.

Hoster berhasil dibuat atau error dikembalikan.
*/
func (s *hosterService) CreateHoster(hoster *model.HosterModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(hoster.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hoster.PasswordHash = string(hash)
	hoster.CreatedAt = time.Now()
	hoster.UpdatedAt = time.Now()

	err = s.repo.CreateHoster(hoster)
	// Cek duplicate email
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgHosterEmailExists)
		}
		return err
	}

	return nil
}

/*
	Mengambil detail hoster dari konteks.

Model hoster dikembalikan jika ditemukan.
*/
func (s *hosterService) GetDetailHoster(ctx context.Context) (*model.HosterModel, error) {
	// Ekstrak ID dari context menggunakan key yang benar (middleware.UserIDKey)
	id, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	hoster, err := s.repo.GetDetailHoster(id)
	if err != nil {
		return nil, err
	}
	if hoster == nil {
		return nil, errors.New("hoster not found")
	}

	return hoster, nil
}

/*
	Struktur untuk respons hoster.

Berisi data token dan informasi pengguna.
*/
type HosterResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
	Antarmuka untuk layanan hoster.

Mendefinisikan metode untuk operasi hoster.
*/
type HosterService interface {
	CreateHoster(*model.HosterModel) error
	LoginHoster(email, password string) (*HosterResponse, error)
	GetDetailHoster(ctx context.Context) (*model.HosterModel, error)
}

/*
	Membuat instance baru dari HosterService.

Instance layanan dikembalikan.
*/
func NewHosterService(repo HosterRepository) HosterService {
	return &hosterService{repo: repo}
}
