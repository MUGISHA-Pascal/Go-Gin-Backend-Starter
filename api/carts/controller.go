package carts

import (
	"fmt"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AddToCart struct {
	ProductId uint `json:"productId"`
	CartId    uint `json:"cartId"`
	Quantity  int  `json:"quantity"`
}
type RemoveItemFromCart struct {
	ProductId uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}

func AddItemToCart(c *gin.Context) {
	var addToCart AddToCart
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	uid, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	if err := c.BindJSON(&addToCart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding the request body"})
		return
	}
	if addToCart.CartId == 0 || addToCart.ProductId == 0 || addToCart.Quantity == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "all details are required"})
		return
	}
	var product database.Product
	if err := database.DB.First(&product, addToCart.ProductId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if product.StockQty < addToCart.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product in stock is not enough"})
		return
	}
	var eCart database.Cart
	if err := database.DB.First(&eCart, addToCart.CartId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			var cart database.Cart
			cart.UserId = uid
			if err := database.DB.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a cart"})
				return
			}
		}
	}
	var existingCartItem database.CartItem
	var cartItem database.CartItem
	if err := database.DB.Where("cart_id = ? and product_id = ?", eCart.ID, addToCart.ProductId).First(&existingCartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cartItem.CartId = eCart.ID
			cartItem.ProductId = addToCart.ProductId
			cartItem.Quantity = addToCart.Quantity
			if err := database.DB.Create(&cartItem).Error; err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error while adding item to cart"})
				return
			}
			var user database.User
			if err := database.DB.First(&user, userId).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			if user.Cart == 0 {
				user.Name = user.Name
				user.Password = user.Password
				user.Email = user.Email
				user.Cart = eCart.ID
				if err := database.DB.Save(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating the user's cart"})
					return
				}
			}
		} else {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking the cart item"})
			return
		}
	} else {
		existingCartItem.Quantity += addToCart.Quantity
		if err := database.DB.Save(&existingCartItem).Error; err != nil {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart item quantity"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "item added successfully"})
}
