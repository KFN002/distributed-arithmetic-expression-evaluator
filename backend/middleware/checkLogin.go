package middleware

import (
	"context"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := models.GetJWTFromSessionStorage(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get JWT token from session storage: %v", err), http.StatusInternalServerError)
			return
		}

		if tokenString == "" {
			ctx := context.WithValue(r.Context(), "userID", 0)
			ctx = context.WithValue(ctx, "login", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, login, err := models.ParseJWT(tokenString)

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse JWT token: %v", err), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "login", login)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
