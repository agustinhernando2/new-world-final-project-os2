// internal/repositories/user_repository.go
package repositories

import (
    "github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).Preload("Orders").First(&user).Error
    return &user, err
}
