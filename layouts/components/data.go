package components

import (
	"fmt"
	"gothstarter/database"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type UserProfileData struct {
	Bio          string
	ProfileImage string
}

type UserDetail struct {
	Location         string
	Activities       []string
	ProvidingService []string
	FindingService   []string
}

type ComponentsData struct {
	isAuthenticated bool `json:"isAuthenticated"`
}

func GetAuth(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		// Handle the case where isAuthenticated is not set
		isAuthenticated = false
	}
	// return ComponentsData{
	// 	isAuthenticated: isAuthenticated,
	// }
	return isAuthenticated
}

func GetUserByCookie(r *http.Request) (*database.User, error) {
	var user_id uint
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		// next.ServeHTTP(w, r.WithContext(ctx))
		// http.Error(w, "There was an error with the cookie when getting it", http.StatusUnauthorized)
		fmt.Printf("Error getting cookie: %v\n", err)
		//		return nil, err
	}
	tokenString := cookie.Value
	claims := &database.UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(*database.UserClaims); ok && token.Valid {
		user_id = claims.UserId
	} else {
		// http.Error(w, "Invalid token", http.StatusUnauthorized)
		fmt.Printf("Invalid token: %v\n", err)
		return nil, err
	}
	query := "SELECT user_id FROM user_tokens WHERE token = $1;"
	if err := database.DB.QueryRow(query, tokenString).Scan(&user_id); err != nil {
		fmt.Printf("Error querying database: %v\n", err)
		return nil, err
	}

	user, err := database.GetUserById(user_id)
	if err != nil {
		// http.Error(w, "user not found", http.StatusNotFound)
		fmt.Printf("User not found: %v\n", err)
		return nil, err
	}
	return user, nil
}
