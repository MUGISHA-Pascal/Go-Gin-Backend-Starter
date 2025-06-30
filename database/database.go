package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=postgres dbname=go-starter port=5432 sslmode=disable"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database " + err.Error())
	}
	DB = connection
	DB.AutoMigrate(&Product{}, &User{}, &Order{}, &OrderItem{}, &Cart{}, &CartItem{}) // to be done after entity creation
}
