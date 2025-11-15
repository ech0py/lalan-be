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

// Struct untuk service admin.
type adminService struct {
	repo AdminRepository
}

// Struct untuk response admin.
type AdminResponse struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Fungsi untuk generate token admin.
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

// Fungsi untuk login admin.
func (s *adminService) LoginAdmin(email, password string) (*AdminResponse, error) {
	admin, err := s.repo.FindByEmailAdminForLogin(email)
	// Cek error atau admin tidak ada
	if err != nil || admin == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verifikasi password
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenAdmin(admin.ID)
}

// Fungsi untuk membuat admin.
func (s *adminService) CreateAdmin(admin *model.AdminModel) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.PasswordHash = string(hash)
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	err = s.repo.CreateAdmin(admin)
	// Cek duplicate email
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgHosterEmailExists)
		}
		return err
	}

	return nil
}

// Fungsi untuk membuat kategori.
func (s *adminService) CreateCategory(category *model.CategoryModel) error {
	existing, err := s.repo.FindCategoryByName(category.Name)
	if err != nil {
		return err
	}
	// Cek nama sudah ada
	if existing != nil {
		return errors.New(message.MsgCategoryNameExists)
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	return s.repo.CreateCategory(category)
}

// Fungsi untuk update kategori.
func (s *adminService) UpdateCategory(category *model.CategoryModel) error {
	category.UpdatedAt = time.Now()
	return s.repo.UpdateCategory(category)
}

// Fungsi untuk hapus kategori.
func (s *adminService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}

// Interface untuk service admin.
type AdminService interface {
	CreateAdmin(*model.AdminModel) error
	LoginAdmin(email, password string) (*AdminResponse, error)
	CreateCategory(*model.CategoryModel) error
	UpdateCategory(*model.CategoryModel) error
	DeleteCategory(id string) error
}

// Fungsi untuk membuat service admin.
func NewAdminService(repo AdminRepository) AdminService {
	return &adminService{repo: repo}
}
