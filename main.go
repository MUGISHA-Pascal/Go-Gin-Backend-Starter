package main

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/MUGISHA-Pascal/Go-Backend-Starter/docs"
	"log"
	"os"
)

// @title           Go Backend Starter API
// @version         1.0
// @description     A RESTful API for e-commerce backend with user authentication, product management, cart operations, and order processing.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error while loading .env file")
	}
	r := gin.Default()
	database.Connect()
	
	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	routes.SetupRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run()
}
