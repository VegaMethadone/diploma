package handlers

import (
	"encoding/json"
	"fmt"
	ownJwt "labyrinth/jwt"
	"labyrinth/logic"
	"net/http"

	"github.com/google/uuid"
)

type RegisterCompanyRequest struct {
	Name        string `json:  "name"`
	Description string `json: "description"`
}

func RegisterCompanyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("jwt_token_user")
	if err != nil {
		http.Error(w, "Missing or invalid JWT cookie", http.StatusUnauthorized)
		return
	}

	claims, err := ownJwt.VerifyToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid or expired JWT token", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req RegisterCompanyRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = logic.NewCompany(req.Name, req.Description, parsedUserID)
	if err != nil {
		ereStr := fmt.Sprintf("Failed to register company: %v", err)
		http.Error(w, ereStr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Company registered successfully"))
}
