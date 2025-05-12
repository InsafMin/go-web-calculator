package auth

import (
	"net/http"
	"strings"

	"github.com/InsafMin/go-web-calculator/pkg/types"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := types.WithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
