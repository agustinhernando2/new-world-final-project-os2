// internal/repositories/item_repository.go
package repositories

import (
	// "fmt"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]*models.Item, error)
	FindByID(itemID uint) (*models.Item, error)
	FindOffersByStatus(status string) ([]models.Item, error)
	UpdateItem(item *models.Item) error
	CreateItem(item *models.Item) error
	GetItemByCategoryAndName(category string, name string) (*models.Item, error)
	ExistsItemByCategoryAndName(category string, name string) (bool, error)
}

type itemRepository struct {
	db *gorm.DB
	// cppserver *cppserver.handler
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}
func (r *itemRepository) FindByID(itemID uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, itemID).Error
	return &item, err
}

func (r *itemRepository) FindOffersByStatus(status string) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("status = ?", status).Find(&items).Error
	return items, err
}

func (r *itemRepository) FindAll() ([]*models.Item, error) {
	var items []*models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *itemRepository) UpdateItemsAvailables() error {
	var items []*models.Item
	err := r.db.Find(&items).Error
	return err
}

func (r *itemRepository) GetItemByCategoryAndName(category string, name string) (*models.Item, error) {
	var item models.Item
	err := r.db.Where("category = ? AND name = ?", category, name).First(&item).Error
	return &item, err
}

func (r *itemRepository) ExistsItemByCategoryAndName(category string, name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Item{}).Where("category = ? AND name = ?", category, name).Count(&count).Error
	return count > 0, err
}

func (r *itemRepository) UpdateItem(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) CreateItem(item *models.Item) error {
	return r.db.Create(item).Error
}
func (r *itemRepository) GetAllItems(db *gorm.DB) ([]models.Item, error) {
	var items []models.Item
	err := db.Model(&models.Item{}).Preload("Items").Find(&items).Error
	return items, err
}
