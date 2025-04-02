package handlers

import (
	"encoding/json"
	"fmt"
	"labyrinth/logic"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json: "email"`
	Password string `json" "password"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	err = logic.NewUser(req.Email, req.Password)
	if err != nil {
		strErr := fmt.Sprintf("Failed to register new user: %v", err)
		http.Error(w, strErr, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}
