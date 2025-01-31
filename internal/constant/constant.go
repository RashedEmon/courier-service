package constant

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// these data should store in database. this is mock data.
var DELIVERY_TYPE_MAPPING = map[int]string{
	1: "SD", // Standard Delivery
	2: "ED", // Express Delivery
	3: "SC", // Scheduled Delivery
	4: "CD", // Cash on Delivery
}

// these data should store in database. this is mock data.
var DELIVERY_STATUS_MAPPING = map[string]string{
	"1": string(StatusPending),
	"2": string(StatusShipped),
	"3": string(StatusDelivered),
	"4": string(StatusCancelled),
}

const (
	OrderTypeDelivery = "Delivery"
	OrderTypePickup   = "Pickup"
	OrderTypeReturn   = "Return"
	OrderTypeExchange = "Exchange"

	ItemTypeParcel      = "Parcel"
	ItemTypeGrocery     = "Grocery"
	ItemTypeElectronics = "Electronics"
)

var ORDER_TYPE_MAPPING = map[string]string{
	"1": OrderTypeDelivery,
	"2": OrderTypePickup,
	"3": OrderTypeReturn,
	"4": OrderTypeExchange,
}

var ITEM_TYPE_MAPPING = map[int]string{
	1: ItemTypeParcel,
	2: ItemTypeGrocery,
	3: ItemTypeElectronics,
}
