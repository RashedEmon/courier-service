package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string          `json:"email" binding:"required,email"`
	Password       string          `json:"password" binding:"required,min=6"`
	DeliveryOrders []DeliveryOrder `gorm:"foreignKey:UserID"`
}
