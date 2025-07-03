package products

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ProductUpdate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockQty    int     `json:"stock_qty"`
}

func CreateProduct(c *gin.Context) {
	var product database.Product
	var eProduct database.Product
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not have access to this feature"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	if user.Role == "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are not allowed for this action"})
		return
	}
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if product.Description == "" || product.Name == "" || product.Price == 0 || product.StockQty == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all products details are required"})
		return
	}
	if err := database.DB.Where("name = ?", product.Name).Find(&eProduct).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting product"})
		return
	}
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving the product!"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product saved successfully", "product": product})
}
func GetAllProducts(c *gin.Context) {
	var products []database.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "products not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "products fetched successfully", "products": products})
}
func GetOneProduct(c *gin.Context) {
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "product id not found"})
		return
	}
	var product database.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product fetched successfully", "product": product})
}
func DeleteProduct(c *gin.Context) {
	productId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id not found"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
		return
	}
	if user.Role != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized for this action"})
	}
	var product database.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "product deleted successfully"})
}
func UpdateProduct(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login first to continue"})
		return
	}
	productId := c.Param("id")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id not provided"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to perform this action"})
		return
	}
	var product database.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	var productUpdateDetails ProductUpdate
	if err := c.ShouldBindJSON(&productUpdateDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if productUpdateDetails.Description != "" {
		product.Description = productUpdateDetails.Description
	}
	if productUpdateDetails.Name != "" {
		product.Name = productUpdateDetails.Name
	}
	if productUpdateDetails.Price != 0 {
		product.Price = productUpdateDetails.Price
	}
	if productUpdateDetails.StockQty != 0 {
		product.StockQty = productUpdateDetails.StockQty
	}
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully", "product": product})
}
