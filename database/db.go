package database

import (
	"fmt"
	"log"
	"os"

	"github.com/pranavpatil6/go_mart/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb(){

	dsn:= os.Getenv("DB_DSN")
	db , err:= gorm.Open(postgres.Open(dsn))
	if err!=nil{
		log.Fatal("Failed to connect to database")
	}
	
	DB = db

	DB.AutoMigrate(&models.User{},&models.Product{},&models.Cart{},&models.CartItem{},&models.Coupon{},models.Order{},models.OrderItem{})
	fmt.Println("connected to db")
}