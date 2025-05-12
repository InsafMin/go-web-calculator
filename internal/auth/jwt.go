package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

func GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": expirationTime.Unix(),
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (int, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id in token")
	}
	return int(userID), nil
}
