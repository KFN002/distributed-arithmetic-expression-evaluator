package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := models.GetJWTFromSessionStorage(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get JWT token from session storage: %v", err), http.StatusInternalServerError)
			return
		}

		if tokenString == "" {
			ctx := context.WithValue(r.Context(), "userID", 0)
			ctx = context.WithValue(ctx, "login", "")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, login, err := models.ParseJWT(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse JWT token: %v", err), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "login", login)

		log.Println(userID, login)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
