package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"courier-service/internal/helpers"
	"courier-service/internal/models"
)

func CreateOrder(c *gin.Context) {
	var order models.DeliveryOrderRequest

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

	createdOrder, err := helpers.StoreOrder(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save data",
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
			"consignment_id":    createdOrder.ConsignmentID,
			"merchant_order_id": createdOrder.MerchantOrderID,
			"order_status":      createdOrder.Status,
			"delivery_fee":      createdOrder.DeliveryFee,
		},
	})

}

func GetOrders(c *gin.Context) {
	//transfer_status=1&archive=0&limit=10&page=2'
	transferStatus := c.DefaultQuery("transfer_status", "1")
	archive := c.DefaultQuery("archive", "0")
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	// function call to get total number of orders
	count := helpers.TotalRowsCount(nil, 0)

	

	orderList := []map[string]string{}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    200,
		"data": gin.H{
			"data":          orderList,
			"total":         4,
			"current_page":  1,
			"per_page":      1,
			"total_in_page": 1,
			"last_page":     4,
		},
	})
}

func CancelOrder(c *gin.Context) {

}
