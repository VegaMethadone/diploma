package auth

import (
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"labyrinth/server/handlers/internal/halper"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func (a AuthHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Check content type
	if err := halper.CheckBodyContent(r); err != nil {
		logger.NewWarnMessage("Invalid content type",
			zap.String("operation", "LoginUserHandler"),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	// 2. Decode request body
	var requestData userLoginRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()
	if err != nil {
		logger.NewWarnMessage("Invalid JSON payload",
			zap.String("operation", "LoginUserHandler"),
			zap.Error(err),
		)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Validate request data
	if err := validateLoginRequest(requestData); err != nil {
		logger.NewWarnMessage("Validation failed",
			zap.String("operation", "LoginUserHandler"),
			zap.Error(err),
			zap.String("email", requestData.Mail),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Process login
	fetchedUser, err := bl.Auth.Login(requestData.Mail, requestData.HashPassword)
	if err != nil {
		logger.NewErrMessage("Login failed",
			zap.String("operation", "LoginUserHandler"),
			zap.Error(err),
			zap.String("email", requestData.Mail),
		)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Generate JWT token
	settings := jwt.MapClaims{
		"id":   fetchedUser.ID,
		"mail": fetchedUser.Email,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":  time.Now().Unix(),
	}
	token := bl.Jwt.NewToken(settings)

	// 6. Set secure HTTP-only cookie
	preparedCookie := &http.Cookie{
		Name:     "labyrinth_user",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
	}
	http.SetCookie(w, preparedCookie)

	// 7. Return success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "success",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "LoginUserHandler"),
			zap.Error(err),
		)
	}

	logger.NewInfoMessage("User logged in successfully",
		zap.String("operation", "LoginUserHandler"),
		zap.String("user_id", fetchedUser.ID.String()),
		zap.String("email", fetchedUser.Email),
	)
}

func validateLoginRequest(req userLoginRequest) error {
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
