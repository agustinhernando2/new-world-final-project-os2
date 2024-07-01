// internal/repositories/order_repository.go
package repositories

import (
	"fmt"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order, oItems []models.OrderItem) error
	FindByID(orderID uint) (*models.Order, error)
	UpdateStatus(orderID uint, status string) error
	GetAllOrders() ([]models.Order, error)
	GetOrders(userID uint) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order, oItems []models.OrderItem) error {
	// Check if oItems is empty
	if len(oItems) == 0 {
		return gorm.ErrEmptySlice
	}

	// Start a transaction
	err := r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&order).Error; err != nil {
			return err // Rollback transaction if order creation fails
		}

		for _, item := range oItems {
			item.OrderID = order.ID
			if err := tx.Create(&item).Error; err != nil {
				return err // Rollback transaction if order item creation fails
			}
		}

		return nil // Commit transaction
	})

	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	// Find the order again to update the order variable
	foundOrder, err := r.FindByID(order.ID)
	if err != nil {
		return fmt.Errorf("error finding order: %w", err)
	}

	// Update the original order variable with the found order
	*order = *foundOrder

	return nil
}

func (r *orderRepository) FindByID(orderID uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.Preload("Items").Joins("JOIN order_items ON order_items.order_id = orders.id").
		Where("orders.id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateStatus(orderID uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Model(&models.Order{}).Preload("Item").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetOrders(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Preload("Item").Find(&orders).Error
	return orders, err
}
