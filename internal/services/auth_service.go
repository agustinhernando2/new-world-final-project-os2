// internal/services/auth_service.go
package services

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(user *models.User) error
	GetOffers() ([]*models.Item, error)
	CheckoutOrders(userID uint, items []models.OrderItem) (*models.Order, error)
	GetOrderStatus(orderID uint) (*models.Order, error)
	GetUserFromId(userID uint) (*models.User, error)
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
	// Copy all fields from u to user
	*user = *u
	return nil
}

func (s *authService) GetOffers() ([]*models.Item, error) {
	return s.itemRepo.FindOffers()
}

func (s *authService) CheckoutOrders(userID uint, oItems []models.OrderItem) (*models.Order, error) {

	// Check if user exists
	_, err := s.userRepo.UserFromId(userID)
	if err != nil {
		return nil, err
	}
	var errItems []models.OrderItem
	for _, item := range oItems {
		i, err := s.itemRepo.FindByID(item.ItemID)
		// If any check fails, add item to errItems
		// Check if items exists
		// Check if Status is "Available"
		// Check if Quantity is greater than 0 and Quantity request in order is less than or equal to 80% of the Quantity in stock
		if err != nil || i.Status != "Available" || i.Quantity <= 0 || item.Quantity <= 0 || float64(item.Quantity) > float64(i.Quantity)*0.8 {
			errItems = append(errItems, item)
			fmt.Println(errItems)
			continue
		}
		// Calculate total price
		item.Price = 10
		item.ItemID = i.ID
	}

	// If there are any errors, return them
	if len(errItems) > 0 {
		// Return error with items that have, errItems json array format
		return nil, fmt.Errorf("error with items: %v", errItems)
	}

	// Create order
	var order = models.Order{
		UserID: userID,
		Status: "Pending",
	}
	// Convertir la estructura a JSON
	orderJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return nil, err
	}
	orderJSON1, err := json.MarshalIndent(oItems, "", "  ")
	if err != nil {
		return nil, err
	}
	// print order in json format using marshal
	fmt.Println("orderJSON", string(orderJSON))
	// print oItems in json format
	fmt.Println("orderJSON1", string(orderJSON1))

	// Call repository to create order
	if err := s.orderRepo.Create(&order, oItems); err != nil {
		return nil, err
	}
	// Update items quantity
	for _, item := range oItems {
		i, _ := s.itemRepo.FindByID(item.ItemID)
		i.Quantity -= item.Quantity
		if err := s.itemRepo.UpdateItem(i); err != nil {
			return nil, err
		}
	}
	return &order, nil
}

func (s *authService) GetOrderStatus(orderID uint) (*models.Order, error) {
	return s.orderRepo.FindByID(orderID)
}

func (s *authService) GetUserFromId(userID uint) (*models.User, error) {
	return s.userRepo.UserFromId(userID)
}
