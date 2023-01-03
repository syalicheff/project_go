package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	ID        uint      `gorm:"primary_key"`
	ProductID uint      `json:"product_id"`
	PricePaid uint      `json:"price_paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
