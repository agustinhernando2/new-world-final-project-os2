// internal/services/auth_service.go
package services

import (
	// "encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/repositories"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(user *models.User) error
	GetOffers() ([]models.Item, error)
	CheckoutOrders(userID uint, items []models.OrderItem) (*models.Order, error)
	GetOrder(userID uint, orderID uint) (*models.Order, error)
	GetOrderStatus(userID uint, orderID uint) (string, error)
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

func (s *authService) GetOffers() ([]models.Item, error) {
	items, err := s.itemRepo.FindOffersByStatus("Available")
	if err != nil {
		return nil, err
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
	// // Convertir la estructura a JSON
	// json, err := json.MarshalIndent(offers, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// // print order in json format using marshal
	// fmt.Println("json", string(json))

	return offers, nil
}

func (s *authService) CheckoutOrders(userID uint, oItems []models.OrderItem) (*models.Order, error) {

	// Check if user exists
	_, err := s.userRepo.UserFromId(userID)
	if err != nil {
		return nil, err
	}
	var errItems []models.OrderItem
	var totalPrice float64
	for index, item := range oItems {
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
		oItems[index].Price = i.Price
		totalPrice += i.Price * float64(item.Quantity)
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
		Total : totalPrice,
	}
	// Convertir la estructura a JSON
	// orderJSON, err := json.MarshalIndent(order, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// print order in json format using marshal
	// fmt.Println("orderJSON", string(orderJSON))

	// Call repository to create order
	if err := s.orderRepo.Create(&order, oItems); err != nil {
		return nil, err
	}
	// Update items quantity
	for _, item := range oItems {
		i, _ := s.itemRepo.FindByID(item.ItemID)
		i.Quantity -= item.Quantity
		// if quantity is 0, change status to "Sold Out"
		if i.Quantity == 0 {
			i.Status = "Sold Out"
		}
		if err := s.itemRepo.UpdateItem(i); err != nil {
			return nil, err
		}
	}
	return &order, nil
}

func (s *authService) GetOrderStatus(userID uint, orderID uint) (string, error) {
	order, err := s.GetOrder(userID, orderID)
	if err != nil {
		return "", err
	}
	return order.Status, nil
}

func (s *authService) GetOrder(userID uint, orderID uint) (*models.Order, error) {
	// Check if user exists
	_, err := s.userRepo.UserFromId(userID)
	if err != nil {
		return nil, err
	}
	// Check if order exists and is owned by user
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not found")
	}
	// // Convertir la estructura a JSON
	// json, err := json.MarshalIndent(order, "", "  ")
	// if err != nil {
	// 	return nil, err
	// }
	// // print order in json format using marshal
	// fmt.Println("json", string(json))
	return order, nil
}

func (s *authService) GetUserFromId(userID uint) (*models.User, error) {
	return s.userRepo.UserFromId(userID)
}
