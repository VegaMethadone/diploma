package middleware

// import (
// 	"context"
// 	ownJwt "labyrinth/jwt"
// 	"net/http"

// 	"github.com/google/uuid"
// )

// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("jwt_token_user")
// 		if err != nil {
// 			http.Error(w, "Missing or invalid JWT cookie", http.StatusUnauthorized)
// 			return
// 		}

// 		claims, err := ownJwt.VerifyToken(cookie.Value)
// 		if err != nil {
// 			http.Error(w, "Invalid or expired JWT token", http.StatusUnauthorized)
// 			return
// 		}

// 		userID, ok := claims["id"].(string)
// 		if !ok {
// 			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
// 			return
// 		}

// 		parsedUUID, err := uuid.Parse(userID)
// 		if err != nil {
// 			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "userUUID", parsedUUID)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
