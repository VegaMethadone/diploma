package handlers

import (
	"fmt"
	"labyrinth/entity/employee"
	"labyrinth/logic"
	"net/http"
)

func GetPositionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	currentEmployee, ok := r.Context().Value("currentEmployee").(employee.Employee)
	if !ok {
		http.Error(w, "Failed to get employee from context", http.StatusBadRequest)
		return
	}

	jsonData, err := logic.GetPositions(&currentEmployee)
	if err != nil {
		cause := fmt.Sprintf("Failed get positions: %v", err)
		http.Error(w, cause, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(jsonData)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
