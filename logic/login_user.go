package logic

import (
	"fmt"

	ownJwt "labyrinth/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func LoginUser(login, password string) (string, error) {
	id, err := ps.LoginUser(login, password)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	claims := jwt.MapClaims{
		"id":    id.String(),
		"login": login,
	}

	token_ := ownJwt.NewToken(claims)
	return token_, nil
}
