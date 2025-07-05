package routes

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/products"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/users"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) *gin.Engine{
	r.POST("/users/login",users.LoginUser)
	r.POST("/users/register",users.RegisterUser)
	protected := r.Group("/")
	protected.Use(middleware.Authentication()){
		setupProductRoutes(protected)
	}
	return r
}
func setupProductRoutes(rg *gin.RouterGroup){
	productRoutes := rg.Group("/products")
	{
	productRoutes.POST("/create",products.CreateProduct)
	productRoutes.GET("/all",products.GetAllProducts)
	productRoutes.GET("/:id",products.GetOneProduct)
	productRoutes.DELETE("/delete/:id",products.DeleteProduct)
	productRoutes.PUT("/update/:id",products.UpdateProduct)
	}
}