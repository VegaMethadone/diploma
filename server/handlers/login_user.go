package handlers

import (
	"encoding/json"
	"fmt"
	"labyrinth/logic"
	"net/http"
	"time"
)

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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

	token, err := logic.LoginUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt_token_user",               // Имя куки
		Value:    token,                          // Значение куки (JWT-токен)
		Expires:  time.Now().Add(24 * time.Hour), // Время жизни куки (24 часа)
		HttpOnly: true,                           // Защита от XSS-атак
		// Secure:   true,                           // Только для HTTPS (если используется)
		Path:     "/",                     // Путь, для которого кука действительна
		SameSite: http.SameSiteStrictMode, // Защита от CSRF-атак
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
