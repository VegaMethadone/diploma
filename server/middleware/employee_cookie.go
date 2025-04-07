package middleware

// import (
// 	"context"
// 	"errors"
// 	"labyrinth/entity/employee"
// 	ownJwt "labyrinth/jwt"
// 	"net/http"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// )

// func AuthMiddlewareEmployee(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("jwt_token_employee")
// 		if err != nil {
// 			http.Error(w, "Missing or invalid JWT cookie", http.StatusUnauthorized)
// 			return
// 		}

// 		claims, err := ownJwt.VerifyToken(cookie.Value)
// 		if err != nil {
// 			http.Error(w, "Invalid or expired JWT token", http.StatusUnauthorized)
// 			return
// 		}
// 		var currentEmployee employee.Employee
// 		err = getEmployeeDataFromCookie(claims, &currentEmployee)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "currentEmployee", currentEmployee)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }

// func getEmployeeDataFromCookie(claims jwt.MapClaims, employeeData *employee.Employee) error {
// 	employeeId, ok := claims["id"].(uuid.UUID)
// 	if !ok {
// 		return errors.New("invalid user ID in token")
// 	}
// 	userId, ok := claims["userId"].(uuid.UUID)
// 	if !ok {
// 		return errors.New("invalid userId in token")
// 	}
// 	companyId, ok := claims["companyId"].(uuid.UUID)
// 	if !ok {
// 		return errors.New("invalid companyId in token")
// 	}
// 	positionId, ok := claims["positionId"].(uuid.UUID)
// 	if !ok {
// 		return errors.New("invalid positionId in token")
// 	}

// 	employeeData.Id = employeeId
// 	employeeData.UserId = userId
// 	employeeData.CompanyId = companyId
// 	employeeData.PositionId = positionId

// 	return nil
// }
