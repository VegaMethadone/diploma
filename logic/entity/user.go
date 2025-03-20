package entity

import (
	"errors"
	"fmt"
	"labyrinth/entiry/user"
	"labyrinth/logic"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewUser(login, passwrod string) (*user.User, error) {
	for _, word := range logic.InjectionKeywords {
		lowerCaseWord := strings.ToLower(word)
		if strings.Contains(login, lowerCaseWord) || strings.Contains(passwrod, lowerCaseWord) {
			err := fmt.Sprintf("login or password should not contain: %s", lowerCaseWord)
			return nil, errors.New(err)
		}
	}
	hashedPassword, err := hashPassword(passwrod)
	if err != nil {
		return nil, err
	}

	return user.NewUser(login, hashedPassword), nil
}
