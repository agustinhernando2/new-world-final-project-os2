// internal/repositories/order_repository.go
package repositories

import (
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
    "gorm.io/gorm"
)

type OrderRepository interface {
    Create(order *models.Order) error
    FindByID(orderID uint) (*models.Order, error)
    UpdateStatus(orderID uint, status string) error
}

type orderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(orderID uint) (*models.Order, error) {
    var order models.Order
    err := r.db.First(&order, orderID).Error
    return &order, err
}

func (r *orderRepository) UpdateStatus(orderID uint, status string) error {
    return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

