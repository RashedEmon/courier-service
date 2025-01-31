package handlers

import (
	"net/http"
	"strconv"

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

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"type":  "error",
			"code":  401,
		},
		)
		return
	}

	userIDInt, _ := userID.(uint)

	createdOrder, err := helpers.StoreOrder(&order, userIDInt)

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
	transferStatus := c.DefaultQuery("transfer_status", "1")
	archive := c.DefaultQuery("archive", "0")
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"type":  "error",
			"code":  401,
		},
		)
		return
	}
	userIDInt, _ := userID.(uint)

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10 // default
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1 // default
	}

	offset := (pageInt - 1) * limitInt

	if archive != "0" {
		archive = "true"
	} else {
		archive = "false"
	}
	// function call to get paginated orders
	orderList, total, err := helpers.GetOrders(transferStatus, archive, limitInt, offset, userIDInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Orders fetched failed",
			"type":    "failed",
			"code":    500,
		})
		return
	}

	currentPage, lastPage := helpers.CalculatePagination(total, offset, limitInt)

	order_data := helpers.PrepareOrderResponse(&orderList)

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    200,
		"data": gin.H{
			"data":          order_data,
			"total":         total,
			"current_page":  currentPage,
			"per_page":      limitInt,
			"total_in_page": len(orderList),
			"last_page":     lastPage,
		},
	})
}

func CancelOrder(c *gin.Context) {
	consignmentID := c.Param("consignment_id")

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"type":  "error",
			"code":  401,
		},
		)
		return
	}

	userIDInt, _ := userID.(uint)

	message, _ := helpers.CancelOder(consignmentID, userIDInt)

	c.JSON(message.Code, gin.H{
		"message": message.Message,
		"type":    message.Type,
		"code":    message.Code,
	},
	)

}
