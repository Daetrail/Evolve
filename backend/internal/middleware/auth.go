package middleware

import (
	"context"
	"evolve-server/internal/util"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// creates a reusable auth middleware function in order to guard multiple routes in the API
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	// actual middleware function, next is the route which needs to be guarded
	return func(next http.Handler) http.Handler {
		// wraps the protected with JWT verification logic
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// fetch token
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				util.SendError(w, http.StatusUnauthorized, "unauthorised")
				return
			}
			
			// verify token
			// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac
			tokenStr := strings.TrimPrefix(header, "Bearer ") 
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err != nil || !token.Valid {
				util.SendError(w, http.StatusUnauthorized, "unauthorised")
				return
			}
			
			// put userID into context
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				util.SendError(w, http.StatusInternalServerError, "internal server error")
				return
			}
			userID := claims["userID"].(string)
			ctx := context.WithValue(r.Context(), "userID", userID)
			
			// Run protected route
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}