package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"courier-service/internal/helpers"
	"courier-service/internal/models"
)

func CreateOrder(c *gin.Context) {
	var order models.DeliveryOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		validationErrors := helpers.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Please fix the given errors",
			"type":    "error",
			"code":    422,
			"errors":  validationErrors,
		})
		return
	}

	if err := helpers.ValidateRequestData(&order); err != nil {
		if validationErr, ok := err.(*helpers.ValidationError); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": validationErr.Message,
				"type":    validationErr.Type,
				"code":    validationErr.Code,
				"errors":  validationErr.Errors,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Validation failed",
				"type":    "error",
				"code":    http.StatusInternalServerError,
			})
		}
		return
	}

	if err := helpers.StoreOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "faield to save data",
			"type":    "error",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    200,
		"data": gin.H{
			"consignment_id":    order.ConsignmentID,
			"merchant_order_id": order.MerchantOrderID,
			"order_status":      "Pending",
			"delivery_fee":      order.DeliveryFee,
		},
	})

}

func GetOrders(c *gin.Context) {

}

func CancelOrder(c *gin.Context) {

}
