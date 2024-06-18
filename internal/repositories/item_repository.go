// internal/repositories/item_repository.go
package repositories

import (
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
    "gorm.io/gorm"
)

type ItemRepository interface {
    FindAll() ([]*models.Item, error)
    FindOffers() ([]*models.Item, error)
}

type itemRepository struct {
    db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
    return &itemRepository{db: db}
}

func (r *itemRepository) FindAll() ([]*models.Item, error) {
    var items []*models.Item
    err := r.db.Find(&items).Error
    return items, err
}

func (r *itemRepository) FindOffers() ([]*models.Item, error) {
    var items []*models.Item
    err := r.db.Where("quantity > 0").Find(&items).Error
    for _, item := range items {
        item.Price *= 0.8 // Aplicar descuento del 20%
    }
    return items, err
}
