package handlers

import (
	"encoding/json"
	"fmt"
	"labyrinth/logic"
	"net/http"

	"github.com/google/uuid"
)

type RegisterCompanyRequest struct {
	Name        string `json:  "name"`
	Description string `json: "description"`
}

func RegisterCompanyHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, ok := r.Context().Value("userUUID").(uuid.UUID)
	if !ok {
		http.Error(w, "Failed to get user UUID", http.StatusInternalServerError)
		return
	}

	var req RegisterCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := logic.NewCompany(req.Name, req.Description, userUUID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to register company: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Company registered successfully"))
}
