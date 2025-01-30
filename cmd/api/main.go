package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"courier-service/config"
	"courier-service/internal/routers"
	"courier-service/internal/database"
)


func main(){
	r := gin.Default()

	// Load config
	config.LoadConfig()

	//initialize database connection
	database.InitializeDB()

	fmt.Println(config.ConfigInstance)
	//Load routes
	routers.GetAppRoutes(r)

	r.Run(":8080")
}
