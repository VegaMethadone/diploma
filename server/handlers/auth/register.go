package auth

import (
	"encoding/json"
	"labyrinth/server/handlers/internal/halper"
	"net/http"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := halper.CheckBodyContent(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	var requestData userRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
	}

}
