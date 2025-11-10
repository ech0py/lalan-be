package service

import (
	"errors"
	"lalan-be/internal/model"
	"lalan-be/internal/repository"
	"lalan-be/pkg/message"
	"strings"

	"github.com/google/uuid"
)

/*
ItemService mendefinisikan operasi untuk layanan item.
*/
type ItemService interface {
	AddItem(input *model.ItemModel) (*model.ItemModel, error)
	GetAllItems() ([]*model.ItemModel, error)
	GetItemByID(id string) (*model.ItemModel, error)
	GetItemsByUserID(userID string) ([]*model.ItemModel, error)
	UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error)
	DeleteItem(id string, userID string) error
}

/*
itemService mengimplementasikan ItemService.
*/
type itemService struct {
	repo repository.ItemRepository
}

/*
NewItemService membuat instance ItemService dengan repository.
*/
func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

/*
AddItem menambahkan item baru.
Mengembalikan model atau error.
*/
func (s *itemService) AddItem(input *model.ItemModel) (*model.ItemModel, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}
	if input.UserID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}
	if input.CategoryID == "" {
		return nil, errors.New(message.MsgCategoryIDRequired)
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

	existing, err := s.repo.FindItemNameByUserID(input.Name, input.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(message.MsgItemNameExists)
	}

	input.ID = uuid.New().String()

	if err := s.repo.CreateItem(input); err != nil {
		return nil, err
	}

	return s.repo.FindByID(input.ID)
}

/*
GetAllItems mendapatkan semua item.
Mengembalikan slice model atau error.
*/
func (s *itemService) GetAllItems() ([]*model.ItemModel, error) {
	return s.repo.FindAll()
}

/*
GetItemByID mendapatkan item berdasarkan ID.
Mengembalikan model atau error.
*/
func (s *itemService) GetItemByID(id string) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}

	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}

	return item, nil
}

/*
GetItemsByUserID mendapatkan semua item berdasarkan user ID.
Mengembalikan slice model atau error.
*/
func (s *itemService) GetItemsByUserID(userID string) ([]*model.ItemModel, error) {
	if userID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}

	return s.repo.FindByUserID(userID)
}

/*
UpdateItem mengupdate item berdasarkan ID dan user ID.
Mengembalikan model atau error.
*/
func (s *itemService) UpdateItem(id string, userID string, input *model.ItemModel) (*model.ItemModel, error) {
	if id == "" {
		return nil, errors.New(message.MsgItemIDRequired)
	}
	if userID == "" {
		return nil, errors.New(message.MsgUserIDRequired)
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)

	if input.Name == "" {
		return nil, errors.New(message.MsgItemNameRequired)
	}
	if input.CategoryID == "" {
		return nil, errors.New(message.MsgCategoryIDRequired)
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

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return nil, errors.New(message.MsgUnauthorized)
	}

	if input.Name != existing.Name {
		duplicate, err := s.repo.FindItemNameByUserID(input.Name, userID)
		if err != nil {
			return nil, err
		}
		if duplicate != nil {
			return nil, errors.New(message.MsgItemNameExists)
		}
	}

	input.ID = id
	input.UserID = userID
	if err := s.repo.Update(input); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

/*
DeleteItem menghapus item berdasarkan ID dan user ID.
Mengembalikan error.
*/
func (s *itemService) DeleteItem(id string, userID string) error {
	if id == "" {
		return errors.New(message.MsgItemIDRequired)
	}
	if userID == "" {
		return errors.New(message.MsgUserIDRequired)
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New(message.MsgItemNotFound)
	}
	if existing.UserID != userID {
		return errors.New(message.MsgUnauthorized)
	}

	return s.repo.Delete(id)
}
