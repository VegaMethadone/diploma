package auth

import (
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/server/handlers/internal/halper"
	"net/http"
	"net/mail"
	"strings"

	"go.uber.org/zap"
)

func (a AuthHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Check content type
	if err := halper.CheckBodyContent(r); err != nil {
		logger.NewWarnMessage("Invalid content type",
			zap.String("operation", "RegisterUserHandler"),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// 2. Decode request body
	var requestData userRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()
	if err != nil {
		logger.NewWarnMessage("Invalid JSON payload",
			zap.String("operation", "RegisterUserHandler"),
			zap.Error(err),
		)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	// 3. Validate request data
	if err := validateRegisterRequest(requestData); err != nil {
		logger.NewWarnMessage("Validation failed",
			zap.String("operation", "RegisterUserHandler"),
			zap.Error(err),
			zap.String("email", requestData.Mail),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Process registration
	// bl := logic.NewBusinessLogic()
	err = bl.Auth.Register(requestData.Mail, requestData.HashPassword, requestData.Phone)
	if err != nil {
		logger.NewErrMessage("Registration failed",
			zap.String("operation", "RegisterUserHandler"),
			zap.Error(err),
			zap.String("email", requestData.Mail),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"status":  "success",
		"message": "User registered successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "RegisterUserHandler"),
			zap.Error(err),
		)
	}

	logger.NewInfoMessage("User registered successfully",
		zap.String("operation", "RegisterUserHandler"),
		zap.String("email", requestData.Mail),
	)
}

// validateRegisterRequest validates registration request data
func validateRegisterRequest(req userRegisterRequest) error {
	if strings.TrimSpace(req.Mail) == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Mail) {
		return errors.New("invalid email format")
	}

	if len(req.HashPassword) == 0 {
		return errors.New("password hash is required")
	}

	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
