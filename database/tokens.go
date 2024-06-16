package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type UserClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type Token struct {
	TokenId   int       `json:"token_id"`
	UserId    uint      `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateTokenString(usr User) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("there was an error loading the secret")
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	expirationTime := time.Now().Add(15 * time.Minute)
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

func MakeToken(tokenString string, usrId uint) (Token, error) {
	// generate token id
	var id int
	err := DB.QueryRow("SELECT COALESCE(MAX(token_id), 0) + 1 FROM user_tokens").Scan(&id)
	if err != nil {
		return Token{}, err
	}

	tokenExpiration := time.Now().Add(15 * time.Minute)

	return Token{
		TokenId:   id,
		UserId:    usrId,
		Token:     tokenString,
		CreatedAt: time.Now(),
		ExpiresAt: tokenExpiration,
	}, err
}

func SaveToken(db *sql.DB, token Token) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE user_id=$1", token.UserId).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user with ID %d does not exist", token.UserId)
	}
	stmt, err := db.Prepare("INSERT INTO user_tokens (user_id, token_id, token, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(token.UserId, token.TokenId, token.Token, token.CreatedAt, token.ExpiresAt)
	if err != nil {
		return err
	}

	return nil
}
