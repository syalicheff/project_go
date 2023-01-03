package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Product struct {
	gorm.Model
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"unique;not null"`
	Price     uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
