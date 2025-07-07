package database

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" example:"iPhone 15"`
	Description string    `json:"description" example:"Latest iPhone model with advanced features"`
	Price       float64   `json:"price" example:"999.99"`
	StockQty    int       `json:"stock_qty" example:"50"`
	CreateAt    time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" gorm:"uniqueIndex" example:"john@example.com"`
	Password  string    `json:"password" example:"hashedpassword"`
	Role      string    `json:"role" gorm:"default:user" example:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      uint      `example:"1"`
}

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	UserId    uint      `json:"user_id" example:"1"`
	Status    string    `json:"status" example:"PENDING"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      uint      `example:"1"`
}
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey" example:"1"`
	OrderId   uint    `json:"order_id" example:"1"`
	ProductId uint    `json:"product_id" example:"1"`
	Quantity  int     `json:"quantity" example:"2"`
	Price     float64 `json:"price" example:"999.99"`
}
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	UserId    uint      `json:"user_id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CartItems []CartItem
}
type CartItem struct {
	ID        uint `json:"id" gorm:"primaryKey" example:"1"`
	CartId    uint `json:"cart_id" example:"1"`
	ProductId uint `json:"product_id" example:"1"`
	Quantity  int  `json:"quantity" example:"2"`
}

type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	OrderID       uint      `json:"order_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
