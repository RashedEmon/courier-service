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
	ID                 uint        `gorm:"primaryKey"`                        // this will index more optimized way & join query will run faster
	ConsignmentID      string      `gorm:"uniqueIndex" json:"consignment_id"` // it will use for lookup
	StoreID            int         `json:"store_id" binding:"required"`
	MerchantOrderID    string      `json:"merchant_order_id"`
	RecipientName      string      `json:"recipient_name" gorm:"size:100" binding:"required"`
	RecipientPhone     string      `json:"recipient_phone" gorm:"size:20" binding:"required"`
	RecipientAddress   string      `json:"recipient_address" gorm:"size:255" binding:"required"`
	RecipientCity      int         `json:"recipient_city"`
	RecipientZone      int         `json:"recipient_zone"`
	RecipientArea      int         `json:"recipient_area"`
	DeliveryType       int         `json:"delivery_type" binding:"required"`
	ItemType           int         `json:"item_type" binding:"required"`
	SpecialInstruction string      `json:"special_instruction" gorm:"size:1000"`
	ItemQuantity       int         `json:"item_quantity" binding:"required,gte=1"`
	ItemWeight         float64     `json:"item_weight" binding:"required,gte=0"`
	AmountToCollect    float64     `json:"amount_to_collect" binding:"required,gt=0"`
	ItemDescription    string      `json:"item_description"`
	DeliveryFee        float64     `json:"delivery_fee" binding:"gte=0"`
	CodFee             float64     `json:"cod_fee" binding:"gte=0"`
	Status             OrderStatus `json:"status" gorm:"size:50"` // enum prevents invalid status
	Discount           float64     `json:"discount" binding:"gte=0"`
	IsArchived         bool        `json:"is_archived"`
	CreatedAt          time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}
