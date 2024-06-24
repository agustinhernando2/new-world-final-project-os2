// internal/services/admin_service.go
package services

import (
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AdminService interface {
	GetDashboard() ([]*models.Order, error)
	UpdateOrderStatus(orderID uint, status string) error
	UpdateStorage() ([]models.Item, error)
}

type adminService struct {
	userRepo  repositories.UserRepository
	itemRepo  repositories.ItemRepository
	orderRepo repositories.OrderRepository
	cppRepo   repositories.CppRepository
}

func NewAdminService(userRepo repositories.UserRepository, itemRepo repositories.ItemRepository, orderRepo repositories.OrderRepository, cppRepo repositories.CppRepository) AdminService {
	return &adminService{
		userRepo:  userRepo,
		itemRepo:  itemRepo,
		orderRepo: orderRepo,
		cppRepo:   cppRepo,
	}
}

func (s *adminService) GetDashboard() ([]*models.Order, error) {
	// Implementar lógica para obtener el estado de todas las órdenes
	return nil, nil
}

func (s *adminService) UpdateOrderStatus(orderID uint, status string) error {
	return s.orderRepo.UpdateStatus(orderID, status)
}

func (s *adminService) UpdateStorage() ([]models.Item, error) {
	// Get supplies from cpp server
	items, err := s.cppRepo.GetSupplies()
	if err != nil {
		return nil, err
	}

	// If an category and item already exists, update the quantity
	for _, item := range items {
		// if item exists, update quantity
		itemExists, err := s.itemRepo.ExistsItemByCategoryAndName(item.Category, item.Name)
		if err != nil {
			return nil, err
		}
		if itemExists {
			// Get item by category and name
			existingItem, err := s.itemRepo.GetItemByCategoryAndName(item.Category, item.Name)
			if err != nil {
				return nil, err
			}
			existingItem.Quantity = existingItem.Quantity + item.Quantity
			err = s.itemRepo.UpdateItem(existingItem) // Update existing item
			if err != nil {
				return nil, err
			}
		} else {
			err := s.itemRepo.CreateItem(&item) // Create new item
			if err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}
