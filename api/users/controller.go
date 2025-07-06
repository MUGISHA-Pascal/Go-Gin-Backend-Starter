package users

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"fmt"
)

// jwtKey will be loaded from environment in each function

type UserUpdate struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"newpassword123"`
	Name     string `json:"name" example:"John Doe"`
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user account with email, name, password, and optional role
// @Tags users
// @Accept json
// @Produce json
// @Param user body database.User true "User registration data"
// @Success 200 {object} map[string]interface{} "User registered successfully with JWT token"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error or email already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/register [post]
func RegisterUser(c *gin.Context) {
	var newUser database.User
	var eUser database.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if newUser.Email == "" || newUser.Name == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or name or password is required"})
		return
	}
	if err := database.DB.Where("email=?", newUser.Email).First(&eUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exist"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking user existence"})
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}
	newUser.Password = string(hashedPass)
	// Allow user to specify their own role during registration
	// If no role is provided, it will default to "user" as defined in the model
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}
	fmt.Println("JWT_SECRET for generation:", jwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    newUser.ID,
		"email": newUser.Email,
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}
	newUser.Password = ""
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": newUser})
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user with email and password, return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body utils.Credentials true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with JWT token"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/login [post]
func LoginUser(c *gin.Context) {
	var cred utils.Credentials
	var user database.User
	if err := c.BindJSON(&cred); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error binding the request credentials"})
		return
	}
	if cred.Email == "" || cred.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all credentials are required"})
		return
	}
	if err := database.DB.Where("email=?", cred.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
		return
	}
	fmt.Println("JWT_SECRET for login generation:", jwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user})
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user details by ID (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UserUpdate true "User update data"
// @Success 200 {object} database.User "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users/update/user/{id} [put]
func UpdateUser(c *gin.Context) {
	var updateUserDetails UserUpdate
	if err := c.ShouldBindJSON(&updateUserDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid json data"})
		return
	}
	userId := c.Param("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not provided"})
		return
	}
	var user database.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if updateUserDetails.Email != "" {
		user.Email = updateUserDetails.Email
	}
	if updateUserDetails.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(updateUserDetails.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing password"})
			return
		}
		user.Password = string(hashedPass)
	}
	if updateUserDetails.Name != "" {
		user.Name = updateUserDetails.Name
	}
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while saving user"})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteYourAccount godoc
// @Summary Delete current user account
// @Description Delete the authenticated user's account
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{} "Account deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /users/delete/myAccount [delete]
func DeleteYourAccount(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting the user account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user account deleted successfully"})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users (admin only)
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized - admin access required"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Security BearerAuth
// @Router /users/all [get]
func GetAllUsers(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	var user database.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are unauthorised to get all the users"})
		return
	}
	var users []database.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"users": users})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users, "message": "users fetched successfully"})
}

// GetYourAccount godoc
// @Summary Get current user account
// @Description Retrieve the authenticated user's account information
// @Tags users
// @Produce json
// @Success 200 {object} database.User "User account information"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /users/mine [get]
func GetYourAccount(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login to access your private information"})
	}
	var user database.User
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
