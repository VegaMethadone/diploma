package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func LoginCompanyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 1. Получаем userUUID из контекста (установленного AuthMiddleware)
	userUUID, ok := r.Context().Value("userUUID").(uuid.UUID)
	if !ok {
		http.Error(w, "Failed to get user UUID from context", http.StatusInternalServerError)
		return
	}

	// // 2. Проверяем принадлежность пользователя к компании
	// company, employee, err := logic.GetUserCompanyInfo(userUUID)
	// if err != nil {
	// 	http.Error(w, "Failed to get company info: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// if company == nil {
	// 	http.Error(w, "User doesn't belong to any company", http.StatusForbidden)
	// 	return
	// }

	// /*
	// 	создания  jwt токена
	// */

	// // 4. Устанавливаем куки для компании
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "jwt_token_company",
	// 	Value:    companyToken,
	// 	Expires:  time.Now().Add(24 * time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteStrictMode,
	// 	Path:     "/",
	// })

	// 5. Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"message":  "Company login successful",
		"company":  company.Name,
		"position": employee.Position,
	})
}
