package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/casnerano/yandex-gophermart/internal/service/token"
)

type ctxUserUUIDType string

const ctxUserUUIDKey ctxUserUUIDType = "user_uuid"

func JWTAuthenticator(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims := token.Claims{}
			jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !jwtToken.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserUUIDKey, claims.UUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserUUID(ctx context.Context) (string, bool) {
	uuid, ok := ctx.Value(ctxUserUUIDKey).(string)
	if !ok {
		return "", false
	}
	return uuid, true
}
