package database

import (
	"courier-service/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() error {
	db := config.ConfigInstance.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable timezone=Asia/Dhaka", db.Host, db.Username, db.Password, db.Name, db.Port)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Connect to database failed")
		return err
	}
	DB = dbConn

	return nil
}
