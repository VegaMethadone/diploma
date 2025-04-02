package handlers

import (
	"labyrinth/logic"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func LoginCompanyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 1. Получаем userUUID из контекста (установленного AuthMiddleware)
	userId, ok := r.Context().Value("userUUID").(uuid.UUID)
	if !ok {
		http.Error(w, "Failed to get user UUID from context", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	companyIdstr, ok := vars["id"]
	if !ok {
		http.Error(w, "Company ID not found in URL", http.StatusBadRequest)
		return
	}

	companyId, err := uuid.Parse(companyIdstr)
	if err != nil {
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return

	}

	token, err := logic.LoginCompany(userId, companyId)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt_token_employee",
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
