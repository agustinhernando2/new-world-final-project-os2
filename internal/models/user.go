package models

import (
	"gorm.io/gorm"
)

// @Description User
type User struct {
	gorm.Model
	ID       uint   `gorm:"unique;primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"text;not null;default:null"`     // @Description Username
	Password string `json:"password" gorm:"size:255;not null;default:null"` // @Description Password
	Email    string `json:"email" gorm:"size:255;not null;default:null"`    // @Description Email
	IsAdmin  bool   `json:"isAdmin" gorm:"not null;default:false"`          // @Description Is Admin, default false
}

func (User) BeforeCreate(db *gorm.DB) error {
	return nil
}
