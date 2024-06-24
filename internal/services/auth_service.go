// internal/services/auth_service.go
package services

import (
	"fmt"
	"regexp"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(user *models.User) error
	GetOffers() ([]*models.Item, error)
	CheckoutOrders(userID uint, orderIDs []uint) error
	GetOrderStatus(orderID uint) (*models.Order, error)
	GetUserFromId(userID uint) (*models.User, error)
	// IsEmailValid(e string) bool
}

type authService struct {
	userRepo  repositories.UserRepository
	itemRepo  repositories.ItemRepository
	orderRepo repositories.OrderRepository
}

func NewAuthService(userRepo repositories.UserRepository, itemRepo repositories.ItemRepository, orderRepo repositories.OrderRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		itemRepo:  itemRepo,
		orderRepo: orderRepo,
	}
}

func isEmailValid(e string) bool {
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(e)
}

func isPasswordValid(p string) bool {
	return regexp.MustCompile(`[a-zA-Z0-9]+`).MatchString(p)
}



func (s *authService) RegisterUser(user *models.User) error {
	// check email format
	if !isEmailValid(user.Email) {
		return fmt.Errorf("invalid email")
	}
	// Check password and username length
	if len(user.Password) < 4 || len(user.Username) < 4 {
		return fmt.Errorf("password or username too short")
	}
	// Check password format
	if !isPasswordValid(user.Password) {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter and one number")
	}

	// Check if email is already registered
	if s.userRepo.IsEmailRegistered(user.Email) {
		return fmt.Errorf("email already registered")
	}

	// Check if user is admin, if not set it to false
	if user.IsAdmin {
		user.IsAdmin = true
	} else {
		user.IsAdmin = false
	}

	// Call repository to create user
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) AuthenticateUser(user *models.User) error {
	// Call repository to create user
	u, err := s.userRepo.UserMatchPassword(user.Email, user.Password)
	if err != nil {
		return err
	}
	user.IsAdmin = u.IsAdmin
	return nil
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

func (s *authService) GetUserFromId(userID uint) (*models.User, error) {
	return s.userRepo.UserFromId(userID)
}
