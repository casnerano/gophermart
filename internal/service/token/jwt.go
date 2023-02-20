package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

func NewJWT(uuid, secret string) (string, error) {
	claims := Claims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
