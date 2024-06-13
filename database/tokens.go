package database

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(usr User) (string, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &UserClaims{
		UserId: usr.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
