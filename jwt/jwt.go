package jwt

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("0987612345574839201")

func NewToken(settings jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, settings)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Failed to create tokenString:", err)
		return ""
	}

	return tokenString
}
