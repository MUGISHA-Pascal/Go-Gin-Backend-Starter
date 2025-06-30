package utils

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var jwtKey = os.Getenv("JWT_SECRET")

func GenerateToken(UserId uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = UserId
	claims["exp"] = time.Now().Add(time.Hour * 24)
	return token.SignedString([]byte(jwtKey))
}
func validateToken(tokenStr string) (uint, error) {}
