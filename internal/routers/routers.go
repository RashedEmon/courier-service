package routers

import (
	"github.com/gin-gonic/gin"

	"courier-service/internal/handlers"
	middleware "courier-service/internal/middlewares"
)

const PATHPREFIX = "/api/v1"

func GetAppRoutes(router *gin.Engine) {

	api := router.Group(PATHPREFIX)
	{

		api.POST("/signup", handlers.SignUp)
		api.POST("/login", handlers.Login)
		// Logout can be implemented using cache server (ex: redis) server by keeping blacklist as jwt can't invalidate token.

		//	Orders routes
		orders := api.Group("/orders")
		{
			orders.Use(middleware.AuthMiddleware())
			orders.POST("", handlers.CreateOrder)
			orders.GET("/all", handlers.GetOrders)
			orders.PUT("/:consignment_id/cancel", handlers.CancelOrder)
		}
	}
}
