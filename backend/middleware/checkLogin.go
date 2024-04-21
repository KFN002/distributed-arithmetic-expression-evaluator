package middleware

import (
	"context"
	"fmt"
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/internal/databaseManager"
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

		// аноним
		if tokenString == "" {
			ctx := context.WithValue(r.Context(), "userID", 0)
			ctx = context.WithValue(ctx, "login", "")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, login, err := models.ParseJWT(tokenString)

		// ошибка или сессия истекла
		if err != nil {
			err := models.ClearJWTSessionStorage(w, r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to clear JWT token: %v", err), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "login", login)

		ok, err := databaseManager.CheckUser(userID, login)

		// пользователь не найден или ошибка поиска
		if ok != true || err != nil {
			err := models.ClearJWTSessionStorage(w, r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to clear JWT token: %v", err), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		log.Println(userID, login)

		// все окей
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
