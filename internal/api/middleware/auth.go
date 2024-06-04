package middleware

import (
	"context"
	"net/http"
	"strings"
	"toy-rental-system/internal/domain/usecase"
	"toy-rental-system/internal/logger"
)

func AuthMiddleware(tokenUsecase usecase.TokenUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := headerParts[1]

			userID, err := tokenUsecase.CheckToken(token, usecase.ScopeAuthentication)
			if err != nil {
				logger.Error("Invalid token:", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
