// internal/services/auth_service.go
package services

import (
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
    "crypto/sha1"
    "fmt"
)

type AuthService interface {
    RegisterUser(email, password string) (*models.User, error)
    AuthenticateUser(email, password string) (*models.User, error)
    GetOffers() ([]*models.Item, error)
    CheckoutOrders(userID uint, orderIDs []uint) error
    GetOrderStatus(orderID uint) (*models.Order, error)
}

type authService struct {
    userRepo repositories.UserRepository
    itemRepo repositories.ItemRepository
    orderRepo repositories.OrderRepository
}

func NewAuthService(userRepo repositories.UserRepository, itemRepo repositories.ItemRepository, orderRepo repositories.OrderRepository) AuthService {
    return &authService{
        userRepo: userRepo,
        itemRepo: itemRepo,
        orderRepo: orderRepo,
    }
}

func (s *authService) RegisterUser(email, password string) (*models.User, error) {
    hashedPassword := fmt.Sprintf("%x", sha1.Sum([]byte(password))) // Ejemplo simple de hash, recomendable usar algo más seguro en producción
    user := &models.User{
        Email: email,
        Password: hashedPassword,
    }
    err := s.userRepo.Create(user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (s *authService) AuthenticateUser(email, password string) (*models.User, error) {
    hashedPassword := fmt.Sprintf("%x", sha1.Sum([]byte(password)))
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, err
    }
    if user.Password != hashedPassword {
        return nil, fmt.Errorf("invalid credentials")
    }
    return user, nil
}

func (s *authService) GetOffers() ([]*models.Item, error) {
    return s.itemRepo.FindOffers()
}

func (s *authService) CheckoutOrders(userID uint, orderIDs []uint) error {
    // Implementar lógica para realizar la compra de órdenes
    return nil
}

func (s *authService) GetOrderStatus(orderID uint) (*models.Order, error) {
    return s.orderRepo.FindByID(orderID)
}
