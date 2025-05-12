package logic

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

type MyJwt struct{}

var secretKey = []byte("0987612345574839201")

func (m MyJwt) NewToken(settings jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, settings)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Failed to create tokenString:", err)
		return ""
	}

	return tokenString
}

func (m MyJwt) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	return claims, nil
}

func NewMyJwt() MyJwt { return MyJwt{} }
