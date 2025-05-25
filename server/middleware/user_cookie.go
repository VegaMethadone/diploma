package middleware

import (
	"context"
	"fmt"
	"labyrinth/logger"
	"labyrinth/logic"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	userIDKey contextKey = "id"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("labyrinth_user")
		if err != nil {
			logger.NewWarnMessage("Missing auth cookie",
				zap.Error(err),
			)
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if cookie.Value == "" {
			logger.NewWarnMessage("Empty cookie value")
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		bl := logic.NewBusinessLogic()
		claims, err := bl.Jwt.VerifyToken(cookie.Value)
		if err != nil {
			logger.NewWarnMessage("Invalid JWT token",
				zap.Error(err),
			)
			http.Error(w, "Invalid or expired authentication token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			logger.NewWarnMessage("Invalid user ID type in token",
				zap.Any("id_type", fmt.Sprintf("%T", claims["id"])),
			)
			http.Error(w, "Invalid user credentials", http.StatusUnauthorized)
			return
		}

		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			logger.NewWarnMessage("Invalid UUID format in token",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			http.Error(w, "Invalid user credentials", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, parsedUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
