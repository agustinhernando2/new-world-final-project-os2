// internal/services/admin_service.go
package services

import (
	"sort"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AdminService interface {
	GetDashboard() ([]models.Item, []models.Order, error)
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

func (s *adminService) GetDashboard() ([]models.Item, []models.Order, error) {
	// get orders
	orders, err := s.orderRepo.GetAllOrders()
	if err != nil {
		return nil, nil, err
	}
	// get offers
	items, err := s.itemRepo.FindOffersByStatus("Available")
	if err != nil {
		return nil, nil, err
	}
	// Modify quantity to show 20%
	for index, item := range items {
		items[index].Quantity = int(float64(item.Quantity) * 0.2) // show 20%
	}
	// Filter items with quantity > 0
	var offers []models.Item
	for _, item := range items {
		if item.Quantity > 0 {
			offers = append(offers, item)
		}
	}
	// order by quantity available

	sort.Slice(offers, func(i, j int) bool {
		return offers[i].Quantity > offers[j].Quantity
	})
	return offers, orders, nil
}

func (s *adminService) UpdateOrderStatus(orderID uint, status string) error {
	// Check if order exists
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	if order.Status == status {
		return nil
	}
	// Id new status is "Deleted", update item quantity
	if status == "Deleted" {
		// Get order items
		orderItems, err := s.orderRepo.GetOrderItems(orderID)
		if err != nil {
			return err
		}
		// Update item quantity
		for _, orderItem := range orderItems {
			item, err := s.itemRepo.FindByID(orderItem.ItemID)
			if err != nil {
				return err
			}
			item.Quantity += orderItem.Quantity
			err = s.itemRepo.UpdateItem(item)
			if err != nil {
				return err
			}
		}
		// delete order
		err = s.orderRepo.DeleteOrder(order)
		if err != nil {
			return err
		}
		return nil
	}
	// Update order status
	err = s.orderRepo.UpdateStatus(orderID, status)
	if err != nil {
		return err
	}
	return nil
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
	// Delete cppserver items
	err = s.cppRepo.DeleteSupplies()
	if err != nil {
		return nil, err
	} 
	return items, nil
}
