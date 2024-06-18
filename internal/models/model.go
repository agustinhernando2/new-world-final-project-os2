package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user record
// @Description Represents a sale record with customer details and total amount
// type User struct {
// 	ID            uint   `gorm:"unique;primaryKey;autoIncrement"`
// 	Username      string `json:"username" gorm:"text;not null;default:null"`       // @Description Customer Email
// 	Password      string `json:"password" gorm:"text;not null;default:null"`       // @Description Customer Email
// 	CustomerEmail string `json:"customer_email" gorm:"text;not null;default:null"` // @Description Customer Email
// 	IsAdmin       bool   `json:"is_admin" gorm:"not null;default:false"`           // @Description Customer Email
// }

type User struct {
	gorm.Model
	ID       uint   `gorm:"unique;primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"text;not null;default:null"`     // @Description Username
	Password string `json:"password" gorm:"size:255;not null;default:null"` // @Description Password
	Email    string `json:"email" gorm:"size:255;default:null"`             // @Description Email
	IsAdmin  bool   `json:"is_admin" gorm:"not null;default:false"`         // @Description Is Admin, default false
}

type Item struct {
	gorm.Model
	ID       uint    `gorm:"unique;primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"text;not null;default:null"`     // @Description Item Name
	Category string  `json:"category" gorm:"text;not null;default:null"` // @Description Item Category
	Price    float64 `json:"price" gorm:"not null;default:0"`            // @Description Item Price
	Quantity int     `json:"quantity" gorm:"not null;default:0"`         // @Description Item Quantity
	Status   string  `json:"status" gorm:"text;not null;default:null"`   // @Description Order Status
	Orders   []Order `json:"orders" gorm:"many2many:order_item;"`        // @Description Orders associated with this item

}

type Order struct {
	gorm.Model
	ID     uint   `gorm:"unique;primaryKey;autoIncrement"`
	UserID uint   `json:"user_id" gorm:"not null;default:0"`        // @Description User ID
	Status string `json:"status" gorm:"text;not null;default:null"` // @Description Order Status
	Items  []Item `json:"items" gorm:"many2many:order_item;"`       // @Description Items associated with this order
}

type OrderItem struct {
	UserID    uint `gorm:"primaryKey"`
	CourseID  uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	Quantity  int `json:"quantity" gorm:"not null;default:0"` // @Description Item Quantity in Order
}

func (Item) BeforeCreate(db *gorm.DB) error {
	return nil
}
