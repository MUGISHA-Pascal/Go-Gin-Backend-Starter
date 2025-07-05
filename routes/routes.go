package routes

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/carts"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/orders"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/products"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/api/users"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) *gin.Engine {
	r.POST("/users/login", users.LoginUser)
	r.POST("/users/register", users.RegisterUser)
	protected := r.Group("/")
	protected.Use(middleware.Authentication())
	{
		setupProductRoutes(protected)
		setupOrderRoutes(protected)
		setupUserRoutes(protected)
		setupCartRoutes(protected)
	}
	return r
}
func setupProductRoutes(rg *gin.RouterGroup) {
	productRoutes := rg.Group("/products")
	{
		productRoutes.POST("/create", products.CreateProduct)
		productRoutes.GET("/all", products.GetAllProducts)
		productRoutes.GET("/:id", products.GetOneProduct)
		productRoutes.DELETE("/delete/:id", products.DeleteProduct)
		productRoutes.PUT("/update/:id", products.UpdateProduct)
	}
}
func setupUserRoutes(rg *gin.RouterGroup) {
	userRoutes := rg.Group("/users")
	{
		userRoutes.GET("/all", users.GetAllUsers)
		userRoutes.GET("/mine", users.GetYourAccount)
		userRoutes.PUT("/update/user/:id", users.UpdateUser)
		userRoutes.DELETE("/delete/myAccount", users.DeleteYourAccount)
	}
}
func setupOrderRoutes(rg *gin.RouterGroup) {
	orderRoutes := rg.Group("/orders")
	{
		orderRoutes.POST("/place-order", orders.PlaceOrder)
		orderRoutes.PUT("/deliver", orders.Deliver)
		orderRoutes.DELETE("/reject", orders.RejectOrder)
	}
}
func setupCartRoutes(rg *gin.RouterGroup) {
	cartRoutes := rg.Group("/carts")
	{
		cartRoutes.POST("/add", carts.AddItemToCart)
		cartRoutes.DELETE("/remove", carts.RemoveItemToCart)
	}
}
