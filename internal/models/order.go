package model

type Order struct {
	ID                 uint    `gorm:"primaryKey"`
	StoreID            int     `json:"store_id"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name" binding:"required"`
	RecipientPhone     string  `json:"recipient_phone" binding:"required"`
	RecipientAddress   string  `json:"recipient_address" binding:"required"`
	RecipientCity      int     `json:"recipient_city"`
	RecipientZone      int     `json:"recipient_zone"`
	RecipientArea      int     `json:"recipient_area"`
	DeliveryType       int     `json:"delivery_type"`
	ItemType           int     `json:"item_type"`
	SpecialInstruction string  `json:"special_instruction"`
	ItemQuantity       int     `json:"item_quantity"`
	ItemWeight         float64 `json:"item_weight"`
	AmountToCollect    float64 `json:"amount_to_collect" binding:"required,gt=0"`
	ItemDescription    string  `json:"item_description"`
	DeliveryFee        float64 `json:"delivery_fee"`
	CodFee             float64 `json:"cod_fee"`
	Status             string  `json:"status"`
	ConsignmentID      string  `json:"consignment_id"`
	Discount           float64 `json:"discount"`
	Archived           bool    `json:"archived"`
	UpdatedAt          string  `json:"updated_at"`
	CreatedAt          string  `json:"created_at"`
}