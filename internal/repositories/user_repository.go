// internal/repositories/user_repository.go
package repositories

import (
    "fmt"

    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
	"golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    IsEmailRegistered(email string) bool
    UserMatchPassword(email, password string) (*models.User, error)
    UserFromId(id uint) (*models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    // Password hashing
	user.Password = hashedPassword(user.Password)
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *userRepository) IsEmailRegistered(email string) bool {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return false
    }
    return true
}

func (r *userRepository) UserMatchPassword(email, password string) (*models.User, error) {
    var user models.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, fmt.Errorf("user not found")
    }
    
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, fmt.Errorf("invalid password")
    }
    return &user, nil
}

func (r *userRepository) UserFromId(id uint) (*models.User, error) {
    var user models.User
    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
        return nil, fmt.Errorf("user not found")
    }
    return &user, nil
}

func hashedPassword(p string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hashedPassword)
}