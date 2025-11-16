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
Struktur ini menyediakan logika bisnis untuk operasi hoster.
*/
type hosterService struct {
	repo HosterRepository
}

/*
Metode untuk menghasilkan token JWT untuk hoster.
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
Metode untuk mengautentikasi hoster dengan email dan password.
Respons token dikembalikan jika berhasil.
*/
func (s *hosterService) LoginHoster(email, password string) (*HosterResponse, error) {
	hoster, err := s.repo.FindByEmailHosterForLogin(email)
	if err != nil || hoster == nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(hoster.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokenHoster(hoster.ID)
}

/*
Metode untuk membuat hoster baru dengan hashing password.
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
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New(message.MsgHosterEmailExists)
		}
		return err
	}

	return nil
}

/*
Metode untuk mengambil detail hoster dari konteks.
Model hoster dikembalikan jika ditemukan.
*/
func (s *hosterService) GetDetailHoster(ctx context.Context) (*model.HosterModel, error) {
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
Metode untuk membuat item baru untuk hoster.
Item divalidasi, dicek duplikasi nama, dan dibuat di database.
*/
func (s *hosterService) CreateItem(ctx context.Context, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}

	if input.CategoryID == "" {
		return nil, errors.New("category ID is required")
	}

	if input.Stock < 0 {
		return nil, errors.New(message.MsgItemStockInvalid)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgItemPricePerDayInvalid)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.MsgItemDepositInvalid)
	}

	existing, err := s.repo.FindItemNameByUserID(input.Name, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgItemNameExists)
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateItem(input); err != nil {
		return nil, err
	}

	return s.repo.FindItemNameByID(input.ID)
}

/*
Metode untuk mengambil item berdasarkan ID.
Model item dikembalikan jika ditemukan.
*/
func (s *hosterService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}

	item, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}

	return item, nil
}

/*
Metode untuk mengambil semua item.
Daftar model item dikembalikan.
*/
func (s *hosterService) GetAllItems() ([]*model.ItemModel, error) {
	return s.repo.GetAllItems()
}

/*
Metode untuk memperbarui item berdasarkan ID.
Item diperbarui jika ditemukan.
*/
func (s *hosterService) UpdateItem(ctx context.Context, id string, input *model.ItemModel) (*model.ItemModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}

	if input.Stock < 0 {
		return nil, errors.New(message.MsgItemStockInvalid)
	}

	if input.PricePerDay < 0 {
		return nil, errors.New(message.MsgItemPricePerDayInvalid)
	}

	if input.Deposit < 0 {
		return nil, errors.New(message.MsgItemDepositInvalid)
	}

	input.ID = id
	input.UserID = userID
	input.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(input); err != nil {
		return nil, err
	}

	return s.repo.FindItemNameByID(id)
}

/*
Metode untuk menghapus item berdasarkan ID.
Item dihapus jika ditemukan dan milik user.
*/
func (s *hosterService) DeleteItem(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New("invalid token claims")
	}

	existing, err := s.repo.FindItemNameByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteItem(id)
}

/*
Metode untuk membuat terms and conditions baru untuk hoster.
Terms and conditions divalidasi dan dibuat di database.
*/
func (s *hosterService) CreateTermsAndConditions(ctx context.Context, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	input.ID = uuid.New().String()
	input.UserID = userID

	if err := s.repo.CreateTermsAndConditions(input); err != nil {
		return nil, err
	}

	return s.repo.FindTermsAndConditionsByID(input.ID)
}

/*
Metode untuk mengambil terms and conditions berdasarkan ID.
Model terms and conditions dikembalikan jika ditemukan.
*/
func (s *hosterService) FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error) {
	if id == "" {
		return nil, errors.New("ID is required")
	}

	tac, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, err
	}
	if tac == nil {
		return nil, errors.New("terms and conditions not found")
	}

	return tac, nil
}

/*
Metode untuk mengambil semua terms and conditions.
Daftar model terms and conditions dikembalikan.
*/
func (s *hosterService) GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error) {
	return s.repo.GetAllTermsAndConditions()
}

/*
Metode untuk memperbarui terms and conditions berdasarkan ID.
Terms and conditions diperbarui jika ditemukan.
*/
func (s *hosterService) UpdateTermsAndConditions(ctx context.Context, id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("terms and conditions not found")
	}
	if existing.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	input.ID = id
	input.UserID = userID

	if err := s.repo.UpdateTermsAndConditions(input); err != nil {
		return nil, err
	}

	return s.repo.FindTermsAndConditionsByID(id)
}

/*
Metode untuk menghapus terms and conditions berdasarkan ID.
Terms and conditions dihapus jika ditemukan dan milik user.
*/
func (s *hosterService) DeleteTermsAndConditions(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return errors.New("invalid token claims")
	}

	existing, err := s.repo.FindTermsAndConditionsByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("terms and conditions not found")
	}
	if existing.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.repo.DeleteTermsAndConditions(id)
}

/*
Struktur untuk respons hoster.
Struktur ini berisi data token dan informasi pengguna.
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
Antarmuka ini mendefinisikan metode untuk operasi hoster.
*/
type HosterService interface {
	CreateHoster(*model.HosterModel) error
	LoginHoster(email, password string) (*HosterResponse, error)
	GetDetailHoster(ctx context.Context) (*model.HosterModel, error)
	CreateItem(ctx context.Context, input *model.ItemModel) (*model.ItemModel, error)
	GetItemByID(id string) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	UpdateItem(ctx context.Context, id string, input *model.ItemModel) (*model.ItemModel, error)
	DeleteItem(ctx context.Context, id string) error
	CreateTermsAndConditions(ctx context.Context, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	FindTermsAndConditionsByID(id string) (*model.TermsAndConditionsModel, error)
	GetAllTermsAndConditions() ([]*model.TermsAndConditionsModel, error)
	UpdateTermsAndConditions(ctx context.Context, id string, input *model.TermsAndConditionsModel) (*model.TermsAndConditionsModel, error)
	DeleteTermsAndConditions(ctx context.Context, id string) error
}

/*
Fungsi untuk membuat instance baru dari HosterService.
Instance layanan dikembalikan.
*/
func NewHosterService(repo HosterRepository) HosterService {
	return &hosterService{repo: repo}
}
