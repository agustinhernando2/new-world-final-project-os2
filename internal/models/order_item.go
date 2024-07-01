package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID       uint    `gorm:"unique;primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"text;not null;default:null"`     // @Description Item Name
	Category string  `json:"category" gorm:"text;not null;default:null"` // @Description Item Category
	Price    float64 `json:"price" gorm:"not null;default:0"`            // @Description Item Price
	Quantity int     `json:"quantity" gorm:"not null;default:0"`         // @Description Item Quantity
	Status   string  `json:"status" gorm:"text;not null;default:null"`   // @Description Order Status
	Orders   []Order `json:"orders" gorm:"many2many:order_items;"`        // @Description Orders associated with this item

}

type Order struct {
	gorm.Model
	ID     uint   `gorm:"unique;primaryKey;autoIncrement"`
	UserID uint   `json:"userId" gorm:"not null;default:0"`         // @Description User ID
	Status string `json:"status" gorm:"text;not null;default:null"` // @Description Order Status
	Items  []Item `json:"items" gorm:"many2many:order_items;"`       // @Description Items associated with this order
}

type OrderItem struct {
	OrderID   uint `gorm:"primaryKey"`
	ItemID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	Price     float64 `json:"price" gorm:"not null;default:0"`    // @Description Item Price
	Quantity  int     `json:"quantity" gorm:"not null;default:0"` // @Description Item Quantity in Order
}

func (Item) BeforeCreate(db *gorm.DB) error {
	return nil
}
