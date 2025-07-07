package utils

import (
	"fmt"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
)

func ProcessPayment(orderID uint, paymentMethod string) (database.Payment, error) {
	var order database.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		return database.Payment{}, fmt.Errorf("order not found")
	}
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
		PaymentMethod: paymentMethod,
		TransactionID: fmt.Sprintf("TXN-%d", order.ID),
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		return database.Payment{}, fmt.Errorf("failed to record payment")
	}
	order.Status = "PAID"
	if err := database.DB.Save(&order).Error; err != nil {
		return payment, fmt.Errorf("failed to update order status")
	}
	return payment, nil
} 