package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// jwtKey will be loaded from environment in each function

func GenerateToken(UserId uint) (string, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		return "", fmt.Errorf("JWT secret not configured")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = UserId
	claims["exp"] = time.Now().Add(time.Hour * 24)
	return token.SignedString([]byte(jwtKey))
}
func ValidateToken(tokenStr string) (uint, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		return 0, fmt.Errorf("JWT secret not configured")
	}
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			return 0, fmt.Errorf("error parsing user id from claims")
		}
		return uint(userIDFloat), nil
	}
	return 0, fmt.Errorf("invalid token")
}
func ParseToken(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, errors.New("empty token string")
	}
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		return 0, fmt.Errorf("JWT secret not configured")
	}
	fmt.Println("JWT_SECRET length:", len(jwtKey))
	fmt.Println("tokenString", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtKey), nil
	})
	fmt.Println("token", token)
	if err != nil {
		fmt.Println("JWT parse error:", err)
		return 0, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims:", claims)
		if idClaim, ok := claims["id"].(float64); ok {
			fmt.Println("Found id claim:", idClaim)
			return uint(idClaim), nil
		}
		fmt.Println("id claim not found or not float64")
		return 0, errors.New("token does not contain valid id")
	}
	
	fmt.Println("Token not valid or claims not ok")
	return 0, errors.New("invalid token")
}

