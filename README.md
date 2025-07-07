# Go Backend Starter

A RESTful API for e-commerce backend with user authentication, product management, cart operations, and order processing.

## Features

- **User Authentication**: JWT-based authentication with user registration and login
- **Product Management**: CRUD operations for products (admin only)
- **Shopping Cart**: Add and remove items from cart
- **Order Processing**: Place, deliver, and reject orders
- **Role-based Access**: Admin and user roles with different permissions
- **Swagger Documentation**: Complete API documentation

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI 3.0
- **Environment**: Environment variables with godotenv

## Getting Started

### Prerequisites

- Go 1.24.1 or higher
- PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-backend-starter
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
Create a `.env` file in the root directory:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret_key
```

4. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Swagger UI

Once the application is running, you can access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

This provides a complete interactive documentation where you can:
- View all available endpoints
- Test API calls directly from the browser
- See request/response schemas
- Authenticate with JWT tokens

### API Endpoints

#### Authentication
- `POST /users/register` - Register a new user
- `POST /users/login` - User login

#### Users (Protected - JWT required)
- `GET /users/mine` - Get current user account
- `PUT /users/update/user/{id}` - Update user (admin only)
- `DELETE /users/delete/myAccount` - Delete current user account
- `GET /users/all` - Get all users (admin only)

#### Products
- `GET /products/all` - Get all products (public)
- `GET /products/{id}` - Get specific product (public)
- `POST /products/create` - Create product (admin only)
- `PUT /products/update/{id}` - Update product (admin only)
- `DELETE /products/delete/{id}` - Delete product (admin only)

#### Cart (Protected - JWT required)
- `POST /carts/add` - Add item to cart
- `DELETE /carts/remove` - Remove item from cart

#### Orders (Protected - JWT required)
- `POST /orders/place-order` - Place new order
- `PUT /orders/deliver` - Deliver order (admin only)
- `DELETE /orders/reject` - Reject order (admin only)
- `POST /orders/pay` - Pay for an order (virtual payment)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. Register or login to get a JWT token
2. Include the token in the Authorization header:
   ```
   Authorization: Bearer <your_jwt_token>
   ```

## Database Models

### User
- `id`: Primary key
- `name`: User's full name
- `email`: Unique email address
- `password`: Hashed password
- `role`: User role (user/admin)
- `cart`: Associated cart ID

### Product
- `id`: Primary key
- `name`: Product name
- `description`: Product description
- `price`: Product price
- `stock_qty`: Available stock quantity

### Cart
- `id`: Primary key
- `user_id`: Associated user ID
- `cart_items`: Array of cart items

### Order
- `id`: Primary key
- `user_id`: Associated user ID
- `status`: Order status (PENDING/DELIVERED)
- `cart`: Associated cart ID

### Payment
- `id`: Primary key
- `order_id`: Associated order ID
- `amount`: Payment amount
- `status`: Payment status (e.g., PAID)
- `payment_method`: Payment method (e.g., virtual_card)
- `transaction_id`: Transaction reference
- `created_at`, `updated_at`: Timestamps

## Development

### Regenerating Swagger Documentation

If you modify the API endpoints or add new ones, regenerate the Swagger documentation:

```bash
~/go/bin/swag init
```

This will update the documentation files in the `docs/` directory.

### Project Structure

```
├── main.go              # Application entry point
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── docs/                # Generated Swagger documentation
├── api/                 # API controllers
│   ├── users/           # User management
│   ├── products/        # Product management
│   ├── carts/           # Cart operations
│   └── orders/          # Order processing
├── database/            # Database models and connection
├── middleware/          # HTTP middleware
├── routes/              # Route definitions
└── utils/               # Utility functions
```

## License

This project is licensed under the Apache 2.0 License. 