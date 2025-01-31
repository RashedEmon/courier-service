package helpers

import (
	"courier-service/internal/constant"
	"courier-service/internal/database"
	"courier-service/internal/models"
	"courier-service/internal/utils"
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ValidationError struct {
	error
	Message string              `json:"message"`
	Type    string              `json:"type"`
	Code    int                 `json:"code"`
	Errors  map[string][]string `json:"errors"`
}

func ValidateRequestData(order *models.DeliveryOrderRequest) error {
	errors := make(map[string][]string)

	matched := regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`).MatchString(order.RecipientPhone)
	if !matched {
		errors["recipient_phone"] = append(errors["recipient_phone"], "Invalid format of the phone number.")
	}

	if len(errors) > 0 {
		return &ValidationError{
			Message: "Please fix the given errors",
			Type:    "error",
			Code:    422,
			Errors:  errors,
		}
	}

	return nil
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			param := fieldErr.Param()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters long", field, param)
			case "max":
				errors[field] = fmt.Sprintf("%s must be at most %s characters long", field, param)
			case "len":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s characters long", field, param)
			case "gt":
				errors[field] = fmt.Sprintf("%s must be greater than %s", field, param)
			case "gte":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
			case "oneof":
				errors[field] = fmt.Sprintf("%s must be one of [%s]", field, param)
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}

func CreateConsignmentId(deliveryType int) string {
	currentDate := time.Now()
	formattedDate := currentDate.Format("020106")
	return constant.DELIVERY_TYPE_MAPPING[deliveryType] + formattedDate + utils.GenerateRandomString(6)
}

func CalculateDeliveryFee(cityID int, weightKg float64) float64 {

	var basePrice float64
	if cityID == 1 {
		basePrice = 60
	} else {
		basePrice = 100
	}

	if cityID == 1 {
		if weightKg <= 0.5 {
			return basePrice
		} else if weightKg <= 1.0 {
			return 70
		}

		extraWeight := weightKg - 1.0
		extraCharge := math.Ceil(extraWeight) * 15
		return 70 + extraCharge
	}

	if weightKg <= 0.5 {
		return basePrice
	} else if weightKg <= 1.0 {
		return basePrice + 10
	}

	extraWeight := weightKg - 1.0
	extraCharge := math.Ceil(extraWeight) * 15
	return (basePrice + 10) + extraCharge
}

func StoreOrder(orderReq *models.DeliveryOrderRequest, userId uint) (models.DeliveryOrder, error) {

	order := models.DeliveryOrder{}

	order.ConsignmentID = CreateConsignmentId(orderReq.DeliveryType)
	order.StoreID = orderReq.StoreID
	order.MerchantOrderID = orderReq.MerchantOrderID
	order.RecipientName = orderReq.RecipientName
	order.RecipientPhone = orderReq.RecipientPhone
	order.RecipientAddress = orderReq.RecipientAddress
	order.RecipientCity = orderReq.RecipientCity
	order.RecipientZone = orderReq.RecipientZone
	order.RecipientArea = orderReq.RecipientArea
	order.DeliveryType = orderReq.DeliveryType
	order.ItemType = orderReq.ItemType
	order.SpecialInstruction = orderReq.SpecialInstruction
	order.ItemQuantity = orderReq.ItemQuantity
	order.ItemWeight = orderReq.ItemWeight
	order.AmountToCollect = orderReq.AmountToCollect
	order.ItemDescription = orderReq.ItemDescription
	order.DeliveryFee = CalculateDeliveryFee(orderReq.RecipientCity, orderReq.ItemWeight)
	order.CodFee = orderReq.AmountToCollect * 0.01
	order.Status = constant.StatusPending
	order.Discount = 0
	order.UserID = &userId
	order.IsArchived = false

	// perform database operation
	if err := database.DB.Create(&order).Error; err != nil {
		fmt.Println("failed to insert row")
		return models.DeliveryOrder{}, err
	}
	return order, nil
}

// get list of orders and total orders
func GetOrders(transferStatus string, archive string, limit int, offset int, userID uint) ([]models.DeliveryOrder, int, error) {
	var orders []models.DeliveryOrder
	var total int64

	// initial query
	query := database.DB.Model(&models.DeliveryOrder{})

	// Apply filters
	transferStatus = constant.DELIVERY_STATUS_MAPPING[transferStatus]
	if transferStatus != "" {
		query = query.Where("status = ?", transferStatus)
	}

	if archive != "" {
		query = query.Where("is_archived = ?", archive)
	}

	query = query.Where("user_id = ?", userID)

	// Get total count of records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func CalculatePagination(totalRecords, offset, limit int) (int, int) {
	if limit <= 0 {
		return 1, 1
	}

	currentPage := (offset / limit) + 1
	lastPage := (totalRecords + limit - 1) / limit

	return currentPage, lastPage
}

// PrepareOrderResponse formats a slice of DeliveryOrder objects into a slice of maps
func PrepareOrderResponse(orders *[]models.DeliveryOrder) []map[string]interface{} {
	var response []map[string]interface{}

	for _, order := range *orders {
		resp := gin.H{
			"order_consignment_id": order.ConsignmentID,
			"order_created_at":     order.CreatedAt.Format("2006-01-02 15:04:05"), // Formatting time
			"order_description":    order.ItemDescription,
			"merchant_order_id":    order.MerchantOrderID,
			"recipient_name":       order.RecipientName,
			"recipient_address":    order.RecipientAddress,
			"recipient_phone":      order.RecipientPhone,
			"order_amount":         order.AmountToCollect,
			"total_fee":            order.DeliveryFee + order.CodFee,
			"instruction":          order.SpecialInstruction,
			"order_type_id":        order.DeliveryType,
			"cod_fee":              order.CodFee,
			"promo_discount":       0,
			"discount":             order.Discount,
			"delivery_fee":         order.DeliveryFee,
			"order_status":         order.Status,
			"order_type":           "Delivery",
			"item_type":            constant.ITEM_TYPE_MAPPING[order.ItemType],
		}
		response = append(response, resp)
	}

	return response
}

func GetOrderByID(ConsignmentID string, userId uint) (*models.DeliveryOrder, error) {
	var order models.DeliveryOrder
	tx := database.DB.Where("consignment_id = ? and user_id = ?", ConsignmentID, userId).First(&order)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &order, nil
}

func UpdateOrderStatus(order *models.DeliveryOrder) error {
	tx := database.DB.Save(order)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type Message struct {
	Message string
	Type    string
	Code    int
}

// cancel order if order is in cancelable state
func CancelOder(consignmentID string, userId uint) (*Message, error) {
	order, err := GetOrderByID(consignmentID, userId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &Message{
				Message: "No data found",
				Type:    "error",
				Code:    404,
			}, err
		} else {
			return &Message{
				Message: "Database query failed. Try again",
				Type:    "error",
				Code:    503,
			}, err
		}
	}

	if order.Status != "pending" {
		return &Message{
			Message: "Please contact cx to cancel order",
			Type:    "error",
			Code:    400,
		}, errors.New("order cannot be canceled")
	}

	order.Status = "cancelled"
	err = UpdateOrderStatus(order)
	if err != nil {
		return nil, err
	}

	return &Message{
		Message: "Order cancelled successfully",
		Type:    "success",
		Code:    200,
	}, nil
}
