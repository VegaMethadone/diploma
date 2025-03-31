package handlers

import (
	"fmt"
	"labyrinth/logic"
	"net/http"

	"github.com/google/uuid"
)

func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := r.Context().Value("userUUID").(uuid.UUID)
	if !ok {
		http.Error(w, "Failed to get user UUID from context", http.StatusInternalServerError)
		return
	}

	jsonData, err := logic.GetUserCompanies(userId)
	if err != nil {
		cause := fmt.Sprintf("Failed get user companies: %v", err)
		http.Error(w, cause, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(jsonData)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
