package main

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading .env file")
	}
	r := gin.Default()
	database.Connect()
	routes.SetupRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run()
}
