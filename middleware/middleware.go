package middleware

import (
	"context"
	"fmt"
	"gothstarter/database"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// // This cannot be the first thing the server wants to the when it is loading the default path
		isAuthenticated := false
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			ctx := context.WithValue(r.Context(), "isAuthenticated", isAuthenticated)

			next.ServeHTTP(w, r.WithContext(ctx))
			// http.Error(w, "There was an error with the cookie", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value
		// Parse the token string
		claims := &database.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				ctx := context.WithValue(r.Context(), "isAuthenticated", isAuthenticated)

				next.ServeHTTP(w, r.WithContext(ctx))
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				// fmt.Printf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			ctx := context.WithValue(r.Context(), "isAuthenticated", isAuthenticated)

			next.ServeHTTP(w, r.WithContext(ctx))
			/* http.Error(w, "auth", http.StatusUnauthorized) */
			fmt.Println("Unauthorized")
			return
		}

		if !token.Valid {
			/* http.Error(w, "auth", http.StatusUnauthorized) */
			fmt.Println("Unauthorized")
			ctx := context.WithValue(r.Context(), "isAuthenticated", isAuthenticated)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		isAuthenticated = true

		ctx := context.WithValue(r.Context(), "isAuthenticated", isAuthenticated)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
