package helpers

import (
	"courier-service/internal/constant"
	"courier-service/internal/database"
	"courier-service/internal/models"
	"courier-service/internal/utils"
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	error
	Message string              `json:"message"`
	Type    string              `json:"type"`
	Code    int                 `json:"code"`
	Errors  map[string][]string `json:"errors"`
}

func ValidateRequestData(order *models.DeliveryOrder) error {
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

	// Format the current date to DDMMYY (with the year as the last two digits)
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

func StoreOrder(order *models.DeliveryOrder) error {
	order.ConsignmentID = CreateConsignmentId(order.DeliveryType)
	order.DeliveryFee = CalculateDeliveryFee(order.RecipientCity, order.ItemWeight)
	order.CodFee = order.AmountToCollect * 0.01

	// perform database operation
	if err := database.DB.Create(order).Error; err != nil {
		fmt.Println("failed to isnert row")
		return err
	}
	return nil
}
