package routers

import (
	"courier-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

const relativePath = "/api/v1"

func GetAppRoutes(router *gin.Engine) {

	api := router.Group(relativePath)
	{
		//	Auth routes
		// auth := api.Group("/auth")
		// {
		// 	auth.POST("/signup", controllers.SignUp)
		// 	auth.POST("/signin", controllers.SignIn)
		// }

		//	Orders routes
		orders := api.Group("/orders")
		{
			orders.POST("", handlers.CreateOrder)
			orders.GET("/all", handlers.GetOrders)
			orders.PUT("/:consignment_id/cancel", handlers.CancelOrder)
		}
	}
}
