package users

import (
	"github.com/MUGISHA-Pascal/Go-Backend-Starter/database"
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
	if newUser.Role != ""{
		newUser.Role = "admin"
	}
	if err := database.DB.Create(&newUser).Error ; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
	return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":newUser.ID,
		"email":newUser.Email
	})
	tokenString,err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}
	newUser.Password=""
	c.JSON(http.StatusOK, gin.H{"token": tokenString,"user": newUser})
}
