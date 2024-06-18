// internal/services/admin_service.go
package services

import (
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AdminService interface {
    GetDashboard() ([]*models.Order, error)
    UpdateOrderStatus(orderID uint, status string) error
}

type adminService struct {
    orderRepo repositories.OrderRepository
}

func NewAdminService(orderRepo repositories.OrderRepository) AdminService {
    return &adminService{
        orderRepo: orderRepo,
    }
}

func (s *adminService) GetDashboard() ([]*models.Order, error) {
    // Implementar lógica para obtener el estado de todas las órdenes
    return nil, nil
}

func (s *adminService) UpdateOrderStatus(orderID uint, status string) error {
    return s.orderRepo.UpdateStatus(orderID, status)
}
