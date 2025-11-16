package admin

import (
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
Struktur untuk respons admin.
Struktur ini berisi data token dan informasi admin.
*/
type AdminResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

/*
Struktur untuk layanan admin.
Struktur ini menyediakan logika bisnis untuk operasi admin.
*/
type adminService struct {
	repo AdminRepository
}

/*
Metode untuk menghasilkan token JWT untuk admin.
Respons token dikembalikan jika berhasil.
*/
func (s *adminService) generateTokenAdmin(userID string) (*AdminResponse, error) {
	exp := time.Now().Add(1 * time.Hour)

	claims := middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, err
	}

	return &AdminResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: uuid.New().String(),
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

/*
Metode untuk mengautentikasi admin dengan email dan password.
Respons token dikembalikan jika berhasil.
*/
func (s *adminService) LoginAdmin(email, password string) (*AdminResponse, error) {
	admin, err := s.repo.FindByEmailAdminForLogin(email)
	if err != nil || admin == nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenAdmin(admin.ID)
}

/*
Metode untuk membuat admin baru dengan hashing password.
Admin berhasil dibuat atau error dikembalikan.
*/
func (s *adminService) CreateAdmin(admin *model.AdminModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.PasswordHash = string(hash)
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	err = s.repo.CreateAdmin(admin)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgHosterEmailExists)
		}
		return err
	}

	return nil
}

/*
Metode untuk membuat kategori baru.
Kategori berhasil dibuat atau error dikembalikan.
*/
func (s *adminService) CreateCategory(category *model.CategoryModel) error {
	existing, err := s.repo.FindCategoryByName(category.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New(message.MsgCategoryNameExists)
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	return s.repo.CreateCategory(category)
}

/*
Metode untuk memperbarui kategori.
Kategori berhasil diperbarui atau error dikembalikan.
*/
func (s *adminService) UpdateCategory(category *model.CategoryModel) error {
	category.UpdatedAt = time.Now()
	return s.repo.UpdateCategory(category)
}

/*
Metode untuk menghapus kategori.
Kategori berhasil dihapus atau error dikembalikan.
*/
func (s *adminService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}

/*
Antarmuka untuk layanan admin.
Antarmuka ini mendefinisikan metode untuk operasi admin.
*/
type AdminService interface {
	CreateAdmin(*model.AdminModel) error
	LoginAdmin(email, password string) (*AdminResponse, error)
	CreateCategory(*model.CategoryModel) error
	UpdateCategory(*model.CategoryModel) error
	DeleteCategory(id string) error
}

/*
Fungsi untuk membuat instance baru dari AdminService.
Instance layanan dikembalikan.
*/
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}
