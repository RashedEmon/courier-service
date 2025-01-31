package models

import (
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type DeliveryOrder struct {
	ID                 uint   `gorm:"primaryKey"`  // this will index more optimized way & join query will run faster
	ConsignmentID      string `gorm:"uniqueIndex"` // it will use for lookup
	StoreID            int
	MerchantOrderID    string
	RecipientName      string `gorm:"size:100"`
	RecipientPhone     string `gorm:"size:20"`
	RecipientAddress   string `gorm:"size:255"`
	RecipientCity      int
	RecipientZone      int
	RecipientArea      int
	DeliveryType       int
	ItemType           int
	SpecialInstruction string `gorm:"size:1000"`
	ItemQuantity       int
	ItemWeight         float64
	AmountToCollect    float64
	ItemDescription    string
	DeliveryFee        float64
	CodFee             float64
	Status             OrderStatus `gorm:"size:50"`
	Discount           float64
	IsArchived         bool
	UserID             *uint
	User               User      `gorm:"constraint:OnDelete:SET NULL;"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

type DeliveryOrderRequest struct {
	ConsignmentID      string  `json:"consignment_id"`
	StoreID            int     `json:"store_id" binding:"required"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name" binding:"required"`
	RecipientPhone     string  `json:"recipient_phone" binding:"required"`
	RecipientAddress   string  `json:"recipient_address" binding:"required"`
	RecipientCity      int     `json:"recipient_city"`
	RecipientZone      int     `json:"recipient_zone"`
	RecipientArea      int     `json:"recipient_area"`
	DeliveryType       int     `json:"delivery_type" binding:"required"`
	ItemType           int     `json:"item_type" binding:"required"`
	SpecialInstruction string  `json:"special_instruction"`
	ItemQuantity       int     `json:"item_quantity" binding:"required,gte=1"`
	ItemWeight         float64 `json:"item_weight" binding:"required,gte=0"`
	AmountToCollect    float64 `json:"amount_to_collect" binding:"required,gt=0"`
	ItemDescription    string  `json:"item_description"`
	DeliveryFee        float64 `json:"delivery_fee" binding:"gte=0"`
	CodFee             float64 `json:"cod_fee" binding:"gte=0"`
	Status             string  `json:"status" binding:"required"`
	Discount           float64 `json:"discount" binding:"gte=0"`
	IsArchived         bool    `json:"is_archived"`
}
