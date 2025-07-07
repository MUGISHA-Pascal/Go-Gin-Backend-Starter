package orders

import (
	"fmt"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DeliverDetails struct {
	Order uint `json:"order" example:"1"`
}

// PlaceOrder godoc
// @Summary Place a new order
// @Description Place an order using items from the user's cart
// @Tags orders
// @Produce json
// @Success 200 {object} map[string]interface{} "Order placed successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User, cart, or cart items not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/place-order [post]
func PlaceOrder(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	var user database.User
	var cart database.Cart
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err := database.DB.First(&cart, user.Cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}
	var cartItems []database.CartItem
	if err := database.DB.Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No cart items found"})
		return
	}
	var order database.Order
	order.Status = "PENDING"
	order.UserId = user.ID
	order.Cart = cart.ID
	if err := database.DB.Create(&order).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while saving order"})
		return
	}
	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while cleaning cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully"})
}

// Deliver godoc
// @Summary Deliver an order
// @Description Mark an order as delivered (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Param order body DeliverDetails true "Order delivery details"
// @Success 200 {object} map[string]interface{} "Order delivered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized - admin access required"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/deliver [put]
func Deliver(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorised to perform this action"})
		return
	}
	var deliverDetails DeliverDetails
	if err := c.BindJSON(&deliverDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding the request body"})
		return
	}
	var order database.Order
	if err := database.DB.First(&order, deliverDetails.Order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting the order"})
		return
	}
	order.Status = "DELIVERED"
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating the order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order delivered successfully"})
}

// RejectOrder godoc
// @Summary Reject an order
// @Description Reject and delete an order (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Param order body DeliverDetails true "Order rejection details"
// @Success 200 {object} map[string]interface{} "Order rejected successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized - admin access required"
// @Failure 404 {object} map[string]interface{} "User or order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/reject [delete]
func RejectOrder(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorised to perform this action"})
		return
	}
	var DeleteDetails DeliverDetails
	if err := c.BindJSON(&DeleteDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while bindind data"})
		return
	}
	var order database.Order
	if err := database.DB.Where("id=?", DeleteDetails.Order).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order rejected successfully"})
}

// PayOrder godoc
// @Summary Pay for an order
// @Description Simulate virtual payment for an order
// @Tags orders
// @Accept json
// @Produce json
// @Param payment body struct{OrderID uint `json:"order_id"`; PaymentMethod string `json:"payment_method"`} true "Payment details"
// @Success 200 {object} map[string]interface{} "Payment successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /orders/pay [post]
func PayOrder(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to continue"})
		return
	}
	var req struct {
		OrderID       uint   `json:"order_id"`
		PaymentMethod string `json:"payment_method"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	var order database.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.UserId != userId && userId != 1 { // 1 is super admin, adjust as needed
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized to pay for this order"})
		return
	}
	// Simulate payment
	amount := 0.0
	var cart database.Cart
	if err := database.DB.First(&cart, order.Cart).Error; err == nil {
		var cartItems []database.CartItem
		database.DB.Where("cart_id = ?", cart.ID).Find(&cartItems)
		for _, item := range cartItems {
			var product database.Product
			if err := database.DB.First(&product, item.ProductId).Error; err == nil {
				amount += float64(item.Quantity) * product.Price
			}
		}
	}
	payment := database.Payment{
		OrderID:       order.ID,
		Amount:        amount,
		Status:        "PAID",
		PaymentMethod: req.PaymentMethod,
		TransactionID: fmt.Sprintf("TXN-%d", order.ID),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record payment"})
		return
	}
	order.Status = "PAID"
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful", "payment": payment})
}
