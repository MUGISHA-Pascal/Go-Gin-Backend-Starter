package users

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var jwtKey = os.Getenv("JWT_SECRET")

type UserUpdate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

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
	if newUser.Role != "" {
		newUser.Role = "admin"
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
		return
	}
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
	c.JSON(http.StatusOK, gin.H{user})

}
